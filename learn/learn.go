package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func runCommand(input string) {
	fmt.Printf("[INTERCEPTED CMD]: %q\n", input)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", input)
	} else {
		cmd = exec.Command("sh", "-c", input)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Printf("[INTERCEPTED OUTPUT]: done (err=%v)\n", err)
	} else {
		fmt.Println("[INTERCEPTED OUTPUT]: done (success)")
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("🟢 Go shell wrapper started. Type commands (or 'exit' to quit):")

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		if input == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		runCommand(input)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Scanner error: %v\n", err)
		os.Exit(1)
	}
}