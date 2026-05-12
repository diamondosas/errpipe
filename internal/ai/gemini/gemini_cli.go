package gemini

import (
	"fmt"
	"os/exec"
	"time"

	"errpipe/internal/utils/sys"
	"errpipe/internal/utils/window"

	"github.com/go-vgo/robotgo"
)

func GeminiCli(errorMessage string) {
    if _, ok := sys.IsInstalled("gemini"); !ok {
        fmt.Println("Gemini CLI is not installed")
        return
    }

    pids, running := sys.IsRunning("gemini")
    if !running {
        fmt.Println("Gemini CLI is not running. Starting it now...")
        cmd := exec.Command("cmd", "/C", "start", "gemini", "--prompt", errorMessage)
        if err := cmd.Run(); err != nil {
            fmt.Printf("Error starting Gemini CLI: %v\n", err)
        }
        return
    }

    fmt.Printf("Gemini CLI is running (PIDs: %v). Snapping window...\n", pids)

    // Build candidate list: gemini PIDs first, then their parents.
    // Console apps don't own their window — the terminal host (conhost/WT) does.
    seen := map[int32]bool{}
    var candidates []int32
    for _, pid := range pids {
        if !seen[pid] {
            candidates = append(candidates, pid)
            seen[pid] = true
        }
        ppid := sys.GetParentPID(pid)
        if ppid > 0 && !seen[ppid] {
            candidates = append(candidates, ppid)
            seen[ppid] = true
        }
    }

    fmt.Printf("DEBUG: Window candidates (gemini + parents): %v\n", candidates)

    var lastErr error
    for _, targetPID := range candidates {
        lastErr = window.BringWindowToFrontByPID(int(targetPID))
        if lastErr == nil {
            fmt.Printf("DEBUG: Successfully brought PID %d to front. Typing error...\n", targetPID)
            time.Sleep(500 * time.Millisecond)
            robotgo.Type(errorMessage)
            return
        }
        fmt.Printf("DEBUG: Could not bring PID %d to front: %v\n", targetPID, lastErr)
    }

    fmt.Println("Unable to bring any gemini process window to front")
}