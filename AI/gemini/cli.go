package gemini

import (
	"fmt"
	"os/exec"
	"time"

	"errpipe/ai/utils"

	"github.com/go-vgo/robotgo"
	
)

func GeminiCli(errorMessage string) {
	if _, ok := utils.IsInstalled("gemini"); !ok {
		fmt.Println("Gemini CLI is not installed")
		return
	}

	pids, running := utils.IsRunning("gemini")

	if !running {
		fmt.Println("Gemini CLI is not running. Starting it now...")
		cmd := exec.Command("cmd", "/C", "start", "gemini", "--prompt", errorMessage)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error starting Gemini CLI: %v\n", err)
		}
		return
	}

	fmt.Printf("Gemini CLI is running (PIDs: %v). Snapping window...\n", pids)
	
	// Iterate through all found PIDs to find the one with the main window
	var err error
	for _, targetPID := range pids {
		err = utils.BringWindowToFrontByPID(int(targetPID))
		if err == nil {
			break
		}
		fmt.Printf("DEBUG: Could not bring PID %d to front, trying next...\n", targetPID)
	}

	if err != nil{
		fmt.Println("Unable to bring any gemini process window to front")
	}else{
		fmt.Printf("DEBUG: Successfully brought window to front. Typing error message: %q\n", errorMessage)
		time.Sleep(500 * time.Millisecond)
		robotgo.Type(errorMessage)
		// robotgo.KeyTap("enter")
	}

}