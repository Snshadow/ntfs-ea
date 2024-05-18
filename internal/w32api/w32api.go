//go:build windows
// +build windows

package w32api

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// ntdll.dll functions

//sys	ntOpenFile(fileHandle *windows.Handle, accessMask uint32, objectAttributes *windows.OBJECT_ATTRIBUTES, ioStatusBlock *windows.IO_STATUS_BLOCK, sharedAccess uint32, openOptions uint32) (err windows.NTStatus) = ntdll.NtOpenFile
//sys	ntClose(fileHandle windows.Handle) (err windows.NTStatus) = ntdll.NtClose
//sys	ntSetEaFile(fileHandle windows.Handle, ioStatusBlock *windows.IO_STATUS_BLOCK, buffer unsafe.Pointer, length uint32) (err windows.NTStatus) = ntdll.NtSetEaFile
//sys	ntQueryEaFile(fileHandle windows.Handle, ioStatusBlock *windows.IO_STATUS_BLOCK, buffer unsafe.Pointer, length uint32, returnSingleEntry bool, eaList unsafe.Pointer, eaListLength uint32, eaIndex *uint32, restartScan bool) (err windows.NTStatus) = ntdll.NtQueryEaFile
//sys	ntQueryInformationFile(fileHandle windows.Handle, ioStatusBlock *windows.IO_STATUS_BLOCK, fileInformation unsafe.Pointer, length uint32, fileInformationClass int32) (err windows.NTStatus) = ntdll.NtQueryInformationFile

func NtOpenFile(accessMask uint32, objectAttributes *windows.OBJECT_ATTRIBUTES, ioStatusBlock *windows.IO_STATUS_BLOCK, sharedAccess uint32, openOptions uint32) (fileHandle windows.Handle, err error) {
	ntStat := ntOpenFile(&fileHandle, accessMask, objectAttributes, ioStatusBlock, sharedAccess, openOptions)

	if ntStat != windows.STATUS_SUCCESS {
		err = ntStat
	}

	return
}

func NtClose(fileHandle windows.Handle) (err error) {
	ntStat := ntClose(fileHandle)

	if ntStat != windows.STATUS_SUCCESS {
		err = ntStat
	}

	return
}

func NtSetEaFile(fileHandle windows.Handle, ioStatusBlock *windows.IO_STATUS_BLOCK, buffer unsafe.Pointer, length uint32) (err error) {
	ntStat := ntSetEaFile(fileHandle, ioStatusBlock, buffer, length)

	if ntStat != windows.STATUS_SUCCESS {
		err = ntStat
	}

	return
}

func NtQueryEaFile(fileHandle windows.Handle, ioStatusBlock *windows.IO_STATUS_BLOCK, buffer unsafe.Pointer, length uint32, returnSingleEntry bool, eaList unsafe.Pointer, eaListLength uint32, eaIndex *uint32, restartScan bool) (err error) {
	ntStat := ntQueryEaFile(fileHandle, ioStatusBlock, buffer, length, returnSingleEntry, eaList, eaListLength, eaIndex, restartScan)

	if ntStat != windows.STATUS_SUCCESS {
		err = ntStat
	}

	return
}

func NtQueryInformationFile(fileHandle windows.Handle, ioStatusBlock *windows.IO_STATUS_BLOCK, fileInformation unsafe.Pointer, length uint32, fileInformationClass int32) (err error) {
	ntStat := ntQueryInformationFile(fileHandle, ioStatusBlock, fileInformation, length, fileInformationClass)

	if ntStat != windows.STATUS_SUCCESS {
		err = ntStat
	}

	return
}
