// Code generated by 'go generate'; DO NOT EDIT.

package w32api

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var (
	modntdll = windows.NewLazySystemDLL("ntdll.dll")

	procNtClose                = modntdll.NewProc("NtClose")
	procNtOpenFile             = modntdll.NewProc("NtOpenFile")
	procNtQueryEaFile          = modntdll.NewProc("NtQueryEaFile")
	procNtQueryInformationFile = modntdll.NewProc("NtQueryInformationFile")
	procNtSetEaFile            = modntdll.NewProc("NtSetEaFile")
)

func ntClose(fileHandle windows.Handle) (err windows.NTStatus) {
	r0, _, _ := syscall.Syscall(procNtClose.Addr(), 1, uintptr(fileHandle), 0, 0)
	err = windows.NTStatus(r0)
	return
}

func ntOpenFile(fileHandle *windows.Handle, accessMask uint32, objectAttributes *windows.OBJECT_ATTRIBUTES, ioStatusBlock *windows.IO_STATUS_BLOCK, sharedAccess uint32, openOptions uint32) (err windows.NTStatus) {
	r0, _, _ := syscall.Syscall6(procNtOpenFile.Addr(), 6, uintptr(unsafe.Pointer(fileHandle)), uintptr(accessMask), uintptr(unsafe.Pointer(objectAttributes)), uintptr(unsafe.Pointer(ioStatusBlock)), uintptr(sharedAccess), uintptr(openOptions))
	err = windows.NTStatus(r0)
	return
}

func ntQueryEaFile(fileHandle windows.Handle, ioStatusBlock *windows.IO_STATUS_BLOCK, buffer unsafe.Pointer, length uint32, returnSingleEntry bool, eaList unsafe.Pointer, eaListLength uint32, eaIndex *uint32, restartScan bool) (err windows.NTStatus) {
	var _p0 uint32
	if returnSingleEntry {
		_p0 = 1
	}
	var _p1 uint32
	if restartScan {
		_p1 = 1
	}
	r0, _, _ := syscall.Syscall9(procNtQueryEaFile.Addr(), 9, uintptr(fileHandle), uintptr(unsafe.Pointer(ioStatusBlock)), uintptr(buffer), uintptr(length), uintptr(_p0), uintptr(eaList), uintptr(eaListLength), uintptr(unsafe.Pointer(eaIndex)), uintptr(_p1))
	err = windows.NTStatus(r0)
	return
}

func ntQueryInformationFile(fileHandle windows.Handle, ioStatusBlock *windows.IO_STATUS_BLOCK, fileInformation unsafe.Pointer, length uint32, fileInformationClass int32) (err windows.NTStatus) {
	r0, _, _ := syscall.Syscall6(procNtQueryInformationFile.Addr(), 5, uintptr(fileHandle), uintptr(unsafe.Pointer(ioStatusBlock)), uintptr(fileInformation), uintptr(length), uintptr(fileInformationClass), 0)
	err = windows.NTStatus(r0)
	return
}

func ntSetEaFile(fileHandle windows.Handle, ioStatusBlock *windows.IO_STATUS_BLOCK, buffer unsafe.Pointer, length uint32) (err windows.NTStatus) {
	r0, _, _ := syscall.Syscall6(procNtSetEaFile.Addr(), 4, uintptr(fileHandle), uintptr(unsafe.Pointer(ioStatusBlock)), uintptr(buffer), uintptr(length), 0, 0)
	err = windows.NTStatus(r0)
	return
}