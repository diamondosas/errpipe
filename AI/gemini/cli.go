package gemini

import (
	"fmt"
	"os/exec"
	"time"

	"errpipe/AI/utils"
	"github.com/go-vgo/robotgo"
)

// GeminiCli handles interaction with the Gemini CLI tool.
func GeminiCli(errorMessage string) {
	// Check if Gemini CLI is installed
	if _, ok := utils.IsInstalled("gemini"); !ok {
		fmt.Println("Gemini Cli is not Installed Please use another AI")
		return
	}

	// Check whether gemini CLI is open or not
	if !utils.IsProcessRunning("gemini") {
		// If it is not open we start a terminal with the parameter "gemini --prompt"
		fmt.Println("Gemini CLI is not running. Starting it now...")
		
		// Using 'start' to open in a new window as shown in learn/learn.go
		cmd := exec.Command("cmd", "/C", "start", "gemini", "--prompt", errorMessage)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error starting Gemini CLI: %v\n", err)
		}
	} else {
		// If it is open we type the error into the screen
		fmt.Println("Gemini CLI is already running. Typing the error...")
		  
		// Give it a moment to be ready/focused (as in learn/learn.go)
		time.Sleep(2 * time.Second)
		
		// Type the error message
		robotgo.Type(errorMessage)
		robotgo.KeyTap("enter")
	}
}


