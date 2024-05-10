//go:build windows
// +build windows

package ntfs_ea

import (
	"unicode/utf16"

	"github.com/Snshadow/ntfs-ea/w32api"
)

const (
	NeedEa = w32api.FILE_NEED_EA
)

type EaInfo struct {
	Flag    uint8
	EaName  string
	EaValue []byte
}

func (ea *EaInfo) convertToWin32() (w32api.FILE_FULL_EA_INFORMATION, error) {
	fullEa := w32api.FILE_FULL_EA_INFORMATION{
		Flag: ea.Flag,
	}

	var nameLength uint8

	for _, r := range ea.EaName {
		bLen := len([]byte(string(r)))

		if bLen == 1 { // ASCII character
			fullEa.EaName = append(fullEa.EaName, r)
		} else { // Non-ASCII character
			u16 := utf16.Encode([]rune{r})

		}
	}

	return fullEa, nil // dummy
}

func SetEaFile(dstPath string, eaInfo EaInfo) error {

}
