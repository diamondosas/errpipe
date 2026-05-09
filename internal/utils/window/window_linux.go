//go:build linux

package window

import (
		"fmt"
		"os/exec"
		"strconv"
		"strings"
)

// findWindow locates the X11 window ID for the given PID using wmctrl.
// wmctrl -l -p output columns: <wid> <desktop> <pid> <host> <title...>
func findWindow(pid int) (string, error) {
		fmt.Printf("DEBUG: findWindow searching for PID: %d\n", pid)

		out, err := exec.Command("wmctrl", "-l", "-p").Output()
		if err != nil {
			return "", fmt.Errorf("wmctrl not found or failed — is it installed? (apt install wmctrl): %w", err)
		}

		var candidateWID string
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			fields := strings.Fields(line)
			// Expect at least: wid, desktop, pid, host, title
			if len(fields) < 5 {
				continue
			}

			wpid, err := strconv.Atoi(fields[2])
			if err != nil {
				continue
			}

			if wpid != pid {
				continue
			}

			wid := fields[0]
			fmt.Printf("DEBUG: Found window WID %s for PID %d\n", wid, pid)

			// Prefer a mapped (visible) window; keep any match as a fallback.
			if isWindowMapped(wid) {
				fmt.Printf("DEBUG: Window %s is mapped (visible), selecting immediately\n", wid)
				return wid, nil
			}

			if candidateWID == "" {
				fmt.Printf("DEBUG: Window %s is unmapped, storing as candidate\n", wid)
				candidateWID = wid
			}
		}

		if candidateWID != "" {
			return candidateWID, nil
		}

		return "", fmt.Errorf("no window found for pid %d", pid)
}

// isWindowMapped uses xwininfo to check whether the window is viewable.
// Returns false on any error so callers treat failures as "not visible".
func isWindowMapped(wid string) bool {
		out, err := exec.Command("xwininfo", "-id", wid).Output()
		if err != nil {
			return false
		}
		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Map State:") {
				fmt.Printf("DEBUG: xwininfo Map State line: %s\n", line)
				return strings.Contains(line, "IsViewable")
			}
		}
		return false
}

// setForeground raises and activates the X11 window identified by wid.
// First it removes the hidden/minimised state, then activates (raises + focuses).
func setForeground(wid string) error {
		fmt.Printf("DEBUG: setForeground attempting to focus window %s\n", wid)

		// Un-minimise / un-hide the window if needed (non-fatal if it fails).
		if err := exec.Command("wmctrl", "-i", "-r", wid, "-b", "remove,hidden,shaded").Run(); err != nil {
			fmt.Printf("DEBUG: wmctrl un-minimise returned error (non-fatal): %v\n", err)
		}

		// Activate (raise + focus).
		if err := exec.Command("wmctrl", "-i", "-a", wid).Run(); err != nil {
			return fmt.Errorf("wmctrl activate failed for WID %s: %w", wid, err)
		}

		fmt.Printf("DEBUG: setForeground succeeded for window %s\n", wid)
		return nil
}

// BringWindowToFrontByPID finds the main window belonging to the given PID
// and brings it to the foreground, restoring it if minimised.
//
// Requires wmctrl and xwininfo:
//
//	apt install wmctrl x11-utils
func BringWindowToFrontByPID(pid int) error {
		fmt.Printf("DEBUG: BringWindowToFrontByPID(pid=%d) started\n", pid)

		wid, err := findWindow(pid)
		if err != nil {
			fmt.Printf("DEBUG: BringWindowToFrontByPID(pid=%d) failed: %v\n", pid, err)
			return err
		}

		return setForeground(wid)
}
