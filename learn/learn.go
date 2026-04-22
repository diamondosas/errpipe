package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
)

// func main() {
// 	// 1. Run the command and capture the text
// 	// We run 'cmd /C echo Hello World' to get the output on Windows
// 	cmd := exec.Command("cmd", "/C", "echo Hello World")
// 	out, err := cmd.Output()
// 	if err != nil {
// 		fmt.Println("Error running command:", err)
// 		return
// 	}

// 	// Clean up the output text
// 	capturedText := strings.TrimSpace(string(out))
// 	fmt.Printf("I captured this message: %s\n", capturedText)

// 	// 2. Open a second terminal window
// 	// 'start cmd' tells Windows to launch a new terminal window
// 	exec.Command("cmd", "/C", "start", "cmd").Run()

// 	// 3. Wait for the window to appear and focus
// 	time.Sleep(2 * time.Second)

// 	// 4. Automatically type the text into the new window
// 	message := "echo The first terminal sent this: " + capturedText
// 	robotgo.TypeStr(message)
// 	robotgo.KeyTap("enter")
// }