//go:build windows
// +build windows

package w32api

const (
	anySize = 1

	FILE_NEED_EA = 0x80 // file should be interpreted with Extended Attributes(EA)
)

const (
	FileEaInformation = 7
)

type FILE_EA_INFORMATION struct {
	EaSize uint32
}

// 4 bytes aligned
type FILE_FULL_EA_INFORMATION struct {
	NextEntryOffset uint32
	Flags           uint8
	EaNameLength    uint8
	EaValueLength   uint16
	// 1 byte for ASCII character, 2 or more bytes for non-ASCII character, looks like the supported character follows the active codepage of the computer.., for English users, it might be cp1252, cp850...
	//
	// As like file names in NTFS, the name of EA is case-insensitive and shown with upper characters when queried.
	EaName [anySize]int8 // EaNameLength[int8]
	//_ [1]byte // '\0'

	/* EaValue [EaValueLength]byte */
}

// 4 bytes aligned
type FILE_GET_EA_INFORMATION struct {
	NextEntryOffset uint32
	EaNameLength    uint8
	EaName          []int8 // [EaNameLength]int8
	//_ [1]byte // null terminator
}
