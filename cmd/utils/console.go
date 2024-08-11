package utils

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// IsFromOwnConsole checks if the console was created solely for this process(e.g. called from gui(explorer.exe) instead from console(cmd.exe or powershell.exe))
func IsFromOwnConsole() bool {
	var pid uint32
	ret, _, _ := windows.NewLazySystemDLL("kernel32.dll").NewProc("GetConsoleProcessList").Call(uintptr(unsafe.Pointer(&pid)), 1)

	return ret < 2
}
