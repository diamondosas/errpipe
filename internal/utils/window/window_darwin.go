//go:build darwin

package window

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// findWindow verifies the PID belongs to a process that System Events can see
// and returns it unchanged, since macOS focuses apps by PID not by window handle.
// Returns an error if no GUI process with that PID is found.
func findWindow(pid int) (int, error) {
	fmt.Printf("DEBUG: findWindow searching for PID: %d\n", pid)

	// Ask System Events for every process whose unix id matches.
	script := fmt.Sprintf(`
		tell application "System Events"
			set matchedProcs to every process whose unix id is %d
			if (count of matchedProcs) is 0 then
				return "-1"
			end if
			set p to item 1 of matchedProcs
			return (unix id of p) as text
		end tell
	`, pid)

	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return -1, fmt.Errorf("osascript failed while looking up PID %d: %w", pid, err)
	}

	foundPID, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil || foundPID == -1 {
		return -1, fmt.Errorf("no GUI process found for pid %d", pid)
	}

	fmt.Printf("DEBUG: findWindow confirmed PID %d is a GUI process\n", foundPID)
	return foundPID, nil
}

// setForeground brings the process with the given PID to the foreground using
// two AppleScript steps:
//  1. Un-minimise every minimised window that belongs to the process.
//  2. Set the process as frontmost so macOS raises it.
func setForeground(pid int) error {
	fmt.Printf("DEBUG: setForeground attempting to focus PID %d\n", pid)

	// Step 1 – restore any minimised windows for this process.
	restoreScript := fmt.Sprintf(`
		tell application "System Events"
			set targetProc to first process whose unix id is %d
			tell targetProc
				repeat with w in every window
					if miniaturized of w is true then
						set miniaturized of w to false
					end if
				end repeat
			end tell
		end tell
	`, pid)

	if err := exec.Command("osascript", "-e", restoreScript).Run(); err != nil {
		// Non-fatal: the process may have no minimised windows.
		fmt.Printf("DEBUG: un-minimise step returned error (non-fatal): %v\n", err)
	}

	// Step 2 – set the process as frontmost.
	activateScript := fmt.Sprintf(`
		tell application "System Events"
			set frontmost of (first process whose unix id is %d) to true
		end tell
	`, pid)

	if err := exec.Command("osascript", "-e", activateScript).Run(); err != nil {
		fmt.Printf("DEBUG: setForeground failed for PID %d: %v\n", pid, err)
		return fmt.Errorf("osascript activate failed for pid %d: %w", pid, err)
	}

	fmt.Printf("DEBUG: setForeground succeeded for PID %d\n", pid)
	return nil
}

// BringWindowToFrontByPID finds the GUI process with the given PID and brings
// its windows to the foreground, restoring them if minimised.
//
// Requires "Accessibility" permission granted to the calling binary in:
//
//	System Settings → Privacy & Security → Accessibility
func BringWindowToFrontByPID(pid int) error {
	fmt.Printf("DEBUG: BringWindowToFrontByPID(pid=%d) started\n", pid)

	confirmedPID, err := findWindow(pid)
	if err != nil {
		fmt.Printf("DEBUG: BringWindowToFrontByPID(pid=%d) failed: %v\n", pid, err)
		return err
	}

	return setForeground(confirmedPID)
}