package main

import (
	"errpipe/internal/ai/gemini"
	"errpipe/internal/cli"
	"fmt"
)

func sendtoAI(errormsg string, config cli.Config){
	fmt.Printf("Sending Error to %s (%s)\n", config.Provider, config.Mode)
	
	switch config.Provider {
	case "Gemini":
		if config.Mode == "Inline CLI Mode"{
			fmt.Println("Coming Soon")
		}else if config.Mode == "Gemini CLI Mode" {
			gemini.GeminiCli(errormsg)
		} else {
			fmt.Println("Web Mode for Gemini is not yet implemented.")
		}
	case "Claude":
		fmt.Println("Claude support is coming soon.")
	case "ChatGPT":
		fmt.Println("ChatGPT support is coming soon.")
	default:
		fmt.Printf("Provider %s is not supported.\n", config.Provider)
	}
}