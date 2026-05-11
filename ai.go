package main

import (
	"errpipe/internal/ai/chatgpt"
	"errpipe/internal/ai/claude"
	"errpipe/internal/ai/gemini"
	"errpipe/internal/cli"
	"fmt"


	"github.com/zendev-sh/goai"
)

func sendtoAI(errormsg string, config cli.Config) {
	fmt.Printf("Sending Error to %s (%s)\n", config.Provider, config.Mode)

	// Validate API Key for Inline CLI Mode
	if config.Mode == "Inline CLI Mode" && config.APIKey == "" {
		fmt.Println("Error: API Key is required for this mode. Please run 'errpipe --init' to configure it.")
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
		fmt.Printf("Provider %s is not supported.\n", config.Provider)
	}
}

func handleInline(errormsg string, config cli.Config) {
	var stream *goai.TextStream
	var err error

	switch config.Provider {
	case "Gemini":
		stream, err = gemini.Stream(config.APIKey, errormsg)
	case "Claude":
		stream, err = claude.Stream(config.APIKey, errormsg)
	case "ChatGPT":
		stream, err = chatgpt.Stream(config.APIKey, errormsg)
	default:
		fmt.Println("Provider not supported for Inline Mode")
		return
	}

	if err != nil {
		fmt.Printf("Error initializing AI stream: %v\n", err)
		return
	}

	fmt.Println("\n--- AI Analysis (Streaming) ---")
	for text := range stream.TextStream() {
		fmt.Print(text)
	}
	fmt.Println("\n-------------------------------")

	if err := stream.Err(); err != nil {
		fmt.Printf("\nStream error occurred: %v\n", err)
	}
}
