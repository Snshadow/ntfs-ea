//go:build windows
// +build windows

package ntfs_ea

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/Snshadow/ntfs-ea/w32api"
	"github.com/nyaosorg/go-windows-mbcs"
)

const (
	NeedEa = w32api.FILE_NEED_EA

	fullInfoHeaderSize = 8 // 4 + 1 + 1 + 2
	getInfoHeaderSize  = 5 // 4 + 1
)

var (
	cp = uintptr(windows.GetACP())
)

// EaInfo is a simplified struct of FILE_FULL_EA_INFORMATION, see https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/wdm/ns-wdm-_file_full_ea_information
type EaInfo struct {
	Flags   uint8
	EaName  string
	EaValue []byte
}

func strToEaNameBuffer(s string) ([]int8, error) {
	buf, err := mbcs.Utf8ToAnsi(s, cp)
	if err != nil {
		return nil, err
	}

	eaName := unsafe.Slice((*int8)(unsafe.Pointer(&buf[0])), len(buf))

	return eaName, nil
}

func (ea *EaInfo) convertToFullInfoPtr() (*w32api.FILE_FULL_EA_INFORMATION, uint32, []byte, error) {
	var eaName []int8
	var fullInfoLen uint32

	eaName, err := strToEaNameBuffer(ea.EaName)
	if err != nil {
		return nil, 0, nil, err
	}

	if len(eaName) > 0xff {
		return nil, 0, nil, fmt.Errorf("EA name is too long")
	}
	if len(ea.EaValue) > 0xffff {
		return nil, 0, nil, fmt.Errorf("EA value is too long")
	}

	fullInfoLen = fullInfoHeaderSize + uint32(len(eaName)) + 1 + uint32(len(ea.EaValue)) // add 1 for null terminator
	fullInfoLen += 4 - (fullInfoLen % 4)                                                 // align by 4 bytes, in case of entries are buffered

	buf := make([]byte, fullInfoLen)
	fullEa := (*w32api.FILE_FULL_EA_INFORMATION)(unsafe.Pointer(&buf[0]))

	fullEa.Flags = ea.Flags
	fullEa.EaNameLength = uint8(len(eaName))
	fullEa.EaValueLength = uint16(len(ea.EaValue))

	copy(unsafe.Slice(&fullEa.EaName[0], fullEa.EaNameLength), eaName)
	copy(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(&fullEa.EaName[0]), fullEa.EaNameLength+1)), fullEa.EaValueLength), ea.EaValue)

	return fullEa, fullInfoLen, buf, nil
}

// EaWriteFile writes EA info into the given path by converting the given eaInfo into buffer that can be used by NtSetEaFile.
// Writing EA with no content will remove the EA with the according EaName if exists, do nothing if the file did not have EA with EaName.
func EaWriteFile(dstPath string, eaInfo EaInfo) error {
	var isb windows.IO_STATUS_BLOCK
	var unicodePath windows.NTUnicodeString

	var openOptions uint32 = windows.FILE_RANDOM_ACCESS | windows.FILE_SYNCHRONOUS_IO_NONALERT

	stat, err := os.Stat(dstPath)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		openOptions |= windows.FILE_DIRECTORY_FILE
	} else {
		openOptions |= windows.FILE_NON_DIRECTORY_FILE
	}

	absPath, err := filepath.Abs(dstPath)
	if err != nil {
		return err
	}
	absPath = "\\??\\\\" + absPath // use NT Namespace

	u16ptr, err := windows.UTF16PtrFromString(absPath)
	if err != nil {
		return err
	}

	windows.RtlInitUnicodeString(&unicodePath, u16ptr)

	objAttr := windows.OBJECT_ATTRIBUTES{
		Length:             uint32(unsafe.Sizeof(windows.OBJECT_ATTRIBUTES{})),
		RootDirectory:      0,
		ObjectName:         &unicodePath,
		Attributes:         windows.OBJ_CASE_INSENSITIVE,
		SecurityDescriptor: nil,
		SecurityQoS:        nil,
	}

	fHnd, err := w32api.NtOpenFile(windows.FILE_GENERIC_READ|windows.FILE_GENERIC_WRITE, &objAttr, &isb, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE, openOptions)
	if err != nil {
		return err
	}

	eaBuf, bufLen, buf, err := eaInfo.convertToFullInfoPtr()
	if err != nil {
		log.Println("failed to prepare ea buffer:", err)
		goto EXIT
	}

	err = w32api.NtSetEaFile(fHnd, &isb, unsafe.Pointer(eaBuf), bufLen)
	if err != nil {
		goto EXIT
	}

	runtime.KeepAlive(buf) // make sure eaName and eaValue is valid until NtSetEaFile is executed.

EXIT:
	closeErr := w32api.NtClose(fHnd)
	if closeErr != nil {
		if err != nil {
			return err
		}
		return closeErr
	}

	return err
}

// AddEaFileWithFile adds EA into file in dst using the content of the given file in src with the given name and flags.
func AddEaFileWithFile(dst string, src string, name string, flags uint8) error {
	buf, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	if len(buf) > 0xffff {
		return fmt.Errorf("file content size is too big")
	}

	eaInfo := EaInfo{
		Flags:   flags,
		EaName:  name,
		EaValue: buf,
	}

	err = EaWriteFile(dst, eaInfo)
	if err != nil {
		return err
	}

	return nil
}

// QueryFileEa queries all EAs in the file in given path and return EaInfo slice which has flag, name, and value of EA.
// If queryName is specified, will only query for EAs that have EaName included in queryName.
func QueryFileEa(path string, queryName ...string) ([]EaInfo, error) {
	var isb windows.IO_STATUS_BLOCK
	var unicodePath windows.NTUnicodeString

	var openOptions uint32 = windows.FILE_RANDOM_ACCESS | windows.FILE_SYNCHRONOUS_IO_NONALERT

	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		openOptions |= windows.FILE_DIRECTORY_FILE
	} else {
		openOptions |= windows.FILE_NON_DIRECTORY_FILE
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	absPath = "\\??\\\\" + absPath // use NT Namespace

	u16ptr, err := windows.UTF16PtrFromString(absPath)
	if err != nil {
		return nil, err
	}

	windows.RtlInitUnicodeString(&unicodePath, u16ptr)

	objAttr := windows.OBJECT_ATTRIBUTES{
		Length:             uint32(unsafe.Sizeof(windows.OBJECT_ATTRIBUTES{})),
		RootDirectory:      0,
		ObjectName:         &unicodePath,
		Attributes:         windows.OBJ_CASE_INSENSITIVE,
		SecurityDescriptor: nil,
		SecurityQoS:        nil,
	}

	fHnd, err := w32api.NtOpenFile(windows.FILE_GENERIC_READ|windows.FILE_GENERIC_WRITE, &objAttr, &isb, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE, openOptions)
	if err != nil {
		return nil, err
	}

	var eaSize uint32
	var eaInfoArr []EaInfo
	var eaInfoPtr *w32api.FILE_FULL_EA_INFORMATION
	buf, eaIndex := []byte(nil), uint32(0)

	var eaListPtr unsafe.Pointer
	var eaListBuf []byte
	var eaList []w32api.FILE_GET_EA_INFORMATION
	var eaListLen uint32

	sz := &w32api.FILE_EA_INFORMATION{}
	err = w32api.NtQueryInformationFile(fHnd, &isb, unsafe.Pointer(sz), uint32(unsafe.Sizeof(*sz)), w32api.FileEaInformation)
	if err != nil {
		eaSize = 0xffff // just set it to maximum value
	} else if sz.EaSize == 0 {
		log.Println("file does not have ea")
		goto EXIT
	}
	eaSize = sz.EaSize

	// if queryName is specified, create eaList for querying
	if len(queryName) != 0 {
		for i, name := range queryName {
			eaName, err := strToEaNameBuffer(name)
			if err != nil {
				log.Println("failed to prepare name for query:", err)
				continue
			}

			eaList = append(eaList, w32api.FILE_GET_EA_INFORMATION{
				NextEntryOffset: 0,
				EaNameLength:    uint8(len(eaName)),
				EaName:          eaName,
			})

			if i != 0 {
				nextOffset := getInfoHeaderSize + eaList[i-1].EaNameLength + 1 // add 1 for null terminator
				nextOffset += 4 - (nextOffset % 4)                             // align with 4 bytes
				eaList[i-1].NextEntryOffset = uint32(nextOffset)
			}
		}

		// convert slice info buffer
		for _, ea := range eaList {
			eaLen := getInfoHeaderSize + ea.EaNameLength + 1
			eaLen += 4 - (eaLen % 4)

			eaBuf := make([]byte, eaLen)
			eaPtr := (*w32api.FILE_GET_EA_INFORMATION)(unsafe.Pointer(&eaBuf[0]))

			eaPtr.NextEntryOffset = ea.NextEntryOffset
			eaPtr.EaNameLength = ea.EaNameLength
			copy(unsafe.Slice((*int8)(unsafe.Add(unsafe.Pointer(eaPtr), getInfoHeaderSize)), ea.EaNameLength+1), ea.EaName)

			eaListLen += uint32(eaLen)
			eaListBuf = append(eaListBuf, eaBuf...)
		}

		eaListPtr = unsafe.Pointer(&eaListBuf[0])
	}

	buf = make([]byte, eaSize)
	err = w32api.NtQueryEaFile(fHnd, &isb, unsafe.Pointer(&buf[0]), eaSize, false, eaListPtr, eaListLen, &eaIndex, false)

	eaInfoPtr = (*w32api.FILE_FULL_EA_INFORMATION)(unsafe.Pointer(&buf[0]))
	for {
		eaInfo := EaInfo{
			Flags: eaInfoPtr.Flags,
		}

		nameBuf := unsafe.Slice((*byte)(unsafe.Pointer(&eaInfoPtr.EaName[0])), eaInfoPtr.EaNameLength)
		name, err := mbcs.AnsiToUtf8(nameBuf, cp)
		if err != nil {
			log.Println("failed to get name of EA:", err)
		} else {
			eaInfo.EaName = name
		}

		eaInfo.EaValue = make([]byte, eaInfoPtr.EaValueLength)
		copy(eaInfo.EaValue, unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(&eaInfoPtr.EaName[0]), eaInfoPtr.EaNameLength+1)), eaInfoPtr.EaValueLength))

		eaInfoArr = append(eaInfoArr, eaInfo)
		if eaInfoPtr.NextEntryOffset == 0 {
			break
		}
	}

EXIT:
	closeErr := w32api.NtClose(fHnd)
	if closeErr != nil {
		if err != nil {
			return eaInfoArr, err
		}
		return eaInfoArr, closeErr
	}

	return eaInfoArr, nil
}
