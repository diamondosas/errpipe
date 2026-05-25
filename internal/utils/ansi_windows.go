//go:build windows

package utils

import (
	"syscall"
	"unsafe"
)

const (
	enableVirtualTerminalProcessing = 0x0004
)

// EnableANSI enables virtual terminal processing on Windows to allow ANSI escape codes.
func EnableANSI() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	handle, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil {
		return // Not a console, skip silently
	}

	var mode uint32
	getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
	// Enable VT processing flag and write it back
	setConsoleMode.Call(uintptr(handle), uintptr(mode|enableVirtualTerminalProcessing))
}
