package chatgpt

import (
	"fmt"
	"os/exec"
	"time"

	"errpipe/internal/utils"
	"errpipe/internal/utils/window"

	"github.com/go-vgo/robotgo"
)

func ChatgptCli(errorMessage string) {
	if _, ok := utils.IsInstalled("chatgpt"); !ok {
		fmt.Println("ChatGPT CLI is not installed")
		return
	}

	pids, running := utils.IsRunning("chatgpt")
	if !running {
		fmt.Println("ChatGPT CLI is not running. Starting it now...")
		cmd := exec.Command("cmd", "/C", "start", "chatgpt", "--prompt", errorMessage)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error starting ChatGPT CLI: %v\n", err)
		}
		return
	}

	fmt.Printf("ChatGPT CLI is running (PIDs: %v). Snapping window...\n", pids)

	seen := map[int32]bool{}
	var candidates []int32
	for _, pid := range pids {
		if !seen[pid] {
			candidates = append(candidates, pid)
			seen[pid] = true
		}
		ppid := utils.GetParentPID(pid)
		if ppid > 0 && !seen[ppid] {
			candidates = append(candidates, ppid)
			seen[ppid] = true
		}
	}

	fmt.Printf("DEBUG: Window candidates (chatgpt + parents): %v\n", candidates)

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

	fmt.Println("Unable to bring any chatgpt process window to front")
}
