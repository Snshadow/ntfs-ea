//go:build windows
// +build windows

package w32api

const (
	FILE_NEED_EA = 0x80 // file should be interpreted with Extended Attributes(EA)
)

// 4 bytes aligned
type FILE_FULL_EA_INFORMATION struct {
	NextEntryOffset uint32
	Flags           uint8
	EaNameLength    uint8
	EaValueLength   uint16
	EaName          []int8 // 1 byte for ASCII character, 2 bytes for non-ASCII character
	EaValue         []byte
}

// 4 bytes aligned
type FILE_GET_EA_INFORMATION struct {
	NextEntryOffset uint32
	EaNameLength    uint8
	EaName          []int8
}
