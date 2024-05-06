//go:build windows
// +build windows

package w32api

const (
	FILE_NEED_EA = 0x80

	anySize = 1
)

// 4 bytes aligned
type FILE_FULL_EA_INFORMATION struct {
	NextEntryOffset uint32
	Flags           uint8
	EaNameLength    uint8
	EaValueLength   uint16
	EaName          [anySize]int8
	EaValue         [anySize]byte
}

// 4 bytes aligned
type FILE_GET_EA_INFORMATION struct {
	NextEntryOffset uint32
	EaNameLength    uint8
	EaName          [anySize]int8
}
