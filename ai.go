package main

import (
	"errpipe/internal/ai/chatgpt"
	"errpipe/internal/ai/claude"
	"errpipe/internal/ai/gemini"
	"errpipe/internal/cli"
	"errpipe/internal/utils"
	"fmt"


	"github.com/zendev-sh/goai"
)

func sendtoAI(errormsg string, config cli.Config) {
	// Validate API Key for Inline CLI Mode
	if config.Mode == "Inline CLI Mode" && config.APIKey == "" {
		utils.PrintError("API Key is required for this mode. Please run 'errpipe --init' to configure it.")
		return
	}

	switch config.Provider {
	case "Gemini":
		if config.Mode == "Inline CLI Mode" {
			handleInline(errormsg, config)
		} else if config.Mode == "Gemini CLI Mode" {
			gemini.GeminiCli(errormsg)
		} else {
			gemini.OpenWeb(errormsg)
		}
	case "Claude":
		if config.Mode == "Inline CLI Mode" {
			handleInline(errormsg, config)
		} else if config.Mode == "Claude CLI Mode" {
			claude.ClaudeCli(errormsg)
		} else {
			claude.OpenWeb(errormsg)
		}
	case "ChatGPT":
		if config.Mode == "Inline CLI Mode" {
			handleInline(errormsg, config)
		} else if config.Mode == "ChatGPT CLI Mode" {
			chatgpt.ChatgptCli(errormsg)
		} else {
			chatgpt.OpenWeb(errormsg)
		}
	default:
		utils.PrintError(fmt.Sprintf("Provider %s is not supported.", config.Provider))
	}
}

func handleInline(errormsg string, config cli.Config) {
	var stream *goai.TextStream
	var err error

	spinner := utils.StartSpinner("Sending to AI...")

	switch config.Provider {
	case "Gemini":
		stream, err = gemini.Stream(config.APIKey, errormsg)
	case "Claude":
		stream, err = claude.Stream(config.APIKey, errormsg)
	case "ChatGPT":
		stream, err = chatgpt.Stream(config.APIKey, errormsg)
	default:
		spinner.Stop()
		utils.PrintError("Provider not supported for Inline Mode")
		return
	}

	spinner.Stop()

	if err != nil {
		utils.PrintError(fmt.Sprintf("Error initializing AI stream: %v", err))
		return
	}

	fmt.Printf("\n\n%s%s--- AI Analysis ---%s\n", utils.Fg(51), utils.Bold(), utils.ResetStr())
	utils.StreamWithHighlighting(stream.TextStream())
	fmt.Printf("\n\n%s%s-------------------%s\n\n", utils.Fg(51), utils.Bold(), utils.ResetStr())

	if err := stream.Err(); err != nil {
		utils.PrintError(fmt.Sprintf("Stream error occurred: %v", err))
	}
}
