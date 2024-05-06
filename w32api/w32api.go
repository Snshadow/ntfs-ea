//go:build windows
// +build windows

package w32api

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// ntdll.dll functions

//sys	ntOpenFile(fileHandle *windows.Handle, accessMask uint32, objectAttributes *windows.OBJECT_ATTRIBUTES, ioStatusBlock *windows.IO_STATUS_BLOCK, sharedAccess uint32, openOptions uint32) (err error) = ntdll.NtOpenFile
//sys	ntClose(fileHandle windows.Handle) (err error) = ntdll.NtClose
//sys	ntSetEaFile(fileHandle windows.Handle, ioStatusBlock *windows.IO_STATUS_BLOCK, buffer unsafe.Pointer, length uint32) (err error) = ntdll.NtSetEaFile
//sys	ntQueryEaFile(fileHandle windows.Handle, ioStatusBlock *windows.IO_STATUS_BLOCK, buffer unsafe.Pointer, length uint32, returnSingleEntry bool, eaList unsafe.Pointer, eaListLength uint32, eaIndex *uint32, restartScan bool) (err error) = ntdll.NtQueryEaFile	

func NtOpenFile(accessMask uint32, objectAttributes *windows.OBJECT_ATTRIBUTES, ioStatusBlock *windows.IO_STATUS_BLOCK, sharedAccess uint32, openOptions uint32) (fileHandle windows.Handle, err error) {
	r := ntOpenFile(&fileHandle, accessMask, objectAttributes, ioStatusBlock, sharedAccess, openOptions)

	if ntStat := r.(windows.NTStatus); ntStat != windows.STATUS_SUCCESS {
		err = ntStat
	}

	return
}

func NtClose(fileHandle windows.Handle) (err error) {
	r := ntClose(fileHandle)

	if ntStat := r.(windows.NTStatus); ntStat != windows.STATUS_SUCCESS {
		err = ntStat
	}

	return
}

func NtSetEaFile(fileHandle windows.Handle, ioStatusBlock *windows.IO_STATUS_BLOCK, buffer unsafe.Pointer, length uint32) (err error) {
	r := ntSetEaFile(fileHandle, ioStatusBlock, buffer, length)

	if ntStat := r.(windows.NTStatus); ntStat != windows.STATUS_SUCCESS {
		err = ntStat
	}

	return
}

func NtQueryEaFile(fileHandle windows.Handle, ioStatusBlock *windows.IO_STATUS_BLOCK, buffer unsafe.Pointer, length uint32, returnSingleEntry bool, eaList unsafe.Pointer, eaListLength uint32, eaIndex *uint32, restartScan bool) (err error) {
	r := ntQueryEaFile(fileHandle, ioStatusBlock, buffer, length, returnSingleEntry, eaList, eaListLength, eaIndex, restartScan)

	if ntStat := r.(windows.NTStatus); ntStat != windows.STATUS_SUCCESS {
		err = ntStat
	}

	return
}
