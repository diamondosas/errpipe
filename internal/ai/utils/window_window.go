//go:build windows

package utils

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	//Activates and displays the window.
	//If the window is minimized or maximized,
	//the system restores it to its original size and position
	swRestore = 9
)

var (
	user32                  = syscall.MustLoadDLL("user32.dll")
	procEnumWindows         = user32.MustFindProc("EnumWindows")
	procGetWindowPID = user32.MustFindProc("GetWindowThreadProcessId")
	procIsWindowVisible = user32.MustFindProc("IsWindowVisible")
	procGetWindow = user32.MustFindProc("GetWindow")


	procSetForegroundWindow = user32.MustFindProc("SetForegroundWindow")
	procShowWindow          = user32.MustFindProc("ShowWindow")

)

func isMainWindow(hwnd syscall.Handle) bool {
	isVisible, _, _ := procIsWindowVisible.Call(uintptr(hwnd))
	gwOwner := uintptr(4)
	isOwned, _, _ := procGetWindow.Call(uintptr(hwnd), gwOwner)
	return isVisible != 0 && isOwned == 0
}

func enumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := procEnumWindows.Call(enumFunc, lparam)
	if r1 == 0 {
		if e1 != nil {
			if errno, ok := e1.(syscall.Errno); ok && errno != 0 {
				err = errno
			}
		}
	}
	return
}


func getWindowThreadProcessId(hwnd syscall.Handle) (int, error) {
	var pid uintptr = 0
	_, _, err := procGetWindowPID.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))
	return int(pid), err
}

func findWindow(pid int) (syscall.Handle, error) {
    fmt.Printf("DEBUG: findWindow searching for PID: %d\n", pid)
    var hwnd syscall.Handle

    cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
        wpid, _ := getWindowThreadProcessId(h)
        if wpid == pid {
            isVisible, _, _ := procIsWindowVisible.Call(uintptr(h))
            gwOwner := uintptr(4)
            isOwned, _, _ := procGetWindow.Call(uintptr(h), gwOwner)

            fmt.Printf("DEBUG: Found window HWND %v for PID %d: visible=%v, owned=%v\n",
                h, pid, isVisible != 0, isOwned != 0)

            // Accept non-owned windows even if currently not visible
            // (console windows are often hidden/minimized, not truly invisible)
            if isOwned == 0 {
                fmt.Printf("DEBUG: Candidate main window HWND %v for PID %d\n", h, pid)
                hwnd = h
                if isVisible != 0 {
                    return 0 // Prefer a visible one — stop immediately
                }
                // Keep enumerating in case a visible one exists
            }
        }
        return 1
    })
    enumWindows(cb, 0)

    if hwnd == 0 {
        return 0, fmt.Errorf("no visible main window found for pid %d", pid)
    }
    return hwnd, nil
}

func setForeground(h syscall.Handle) error {
	fmt.Printf("DEBUG: setForeground attempting to show and focus HWND %v\n", h)
	procShowWindow.Call(uintptr(h), swRestore)
	r, _, _ := procSetForegroundWindow.Call(uintptr(h))
	if r == 0 {
		fmt.Printf("DEBUG: SetForegroundWindow failed for HWND %v\n", h)
	} else {
		fmt.Printf("DEBUG: SetForegroundWindow succeeded for HWND %v\n", h)
	}
	return nil
}

func BringWindowToFrontByPID(pid int) error {
	fmt.Printf("DEBUG: BringWindowToFrontByPID(pid=%d) started\n", pid)
	h, err := findWindow(pid)
	if err != nil {
		fmt.Printf("DEBUG: BringWindowToFrontByPID(pid=%d) failed: %v\n", pid, err)
		return err
	}
	return setForeground(h)
}