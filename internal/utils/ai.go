package utils

import (
	"context"
	"errpipe/internal/ai/chatgpt"
	"errpipe/internal/ai/claude"
	"errpipe/internal/ai/free"
	"errpipe/internal/ai/gemini"
	"errpipe/internal/cli"
	"fmt"

	"github.com/zendev-sh/goai"
)

func SendToAI(ctx context.Context, errormsg string, config cli.Config) {
	// Validate API Key for Inline CLI Mode, unless it's Free mode
	if config.Mode == "Inline CLI Mode" && config.APIKey == "" && config.Provider != "Free" {
		PrintError("API Key is required for this mode. Please run 'errpipe --init' to configure it.")
		return
	}

	switch config.Provider {
	case "Free":
		if config.Mode == "Inline CLI Mode" {
			HandleFreeInline(ctx, errormsg)
		} else {
			PrintError("Free mode only supports Inline CLI Mode")
		}
	case "Gemini":
		if config.Mode == "Inline CLI Mode" {
			HandleInline(ctx, errormsg, config)
		} else if config.Mode == "Gemini CLI Mode" {
			gemini.GeminiCli(errormsg)
		} else {
			gemini.OpenWeb(errormsg)
		}
	case "Claude":
		if config.Mode == "Inline CLI Mode" {
			HandleInline(ctx, errormsg, config)
		} else if config.Mode == "Claude CLI Mode" {
			claude.ClaudeCli(errormsg)
		} else {
			claude.OpenWeb(errormsg)
		}
	case "ChatGPT":
		if config.Mode == "Inline CLI Mode" {
			HandleInline(ctx, errormsg, config)
		} else if config.Mode == "ChatGPT CLI Mode" {
			chatgpt.ChatgptCli(errormsg)
		} else {
			chatgpt.OpenWeb(errormsg)
		}
	default:
		PrintError(fmt.Sprintf("Provider %s is not supported.", config.Provider))
	}
}

func HandleFreeInline(ctx context.Context, errormsg string) {
	outChan := make(chan string)
	errChan := make(chan error)

	spinner := StartSpinner()

	go free.StreamToChan(errormsg, outChan, errChan)

	// We stop the spinner before we start printing the output.
	// But we need to wait for the first chunk to stop the spinner ideally.
	// For simplicity, stop it immediately before streaming.
	spinner.Stop()

	fmt.Printf("\n\n%s%s--- AI Analysis ---%s\n", Fg(51), Bold(), ResetStr())

	go func() {
		for err := range errChan {
			if err != nil {
				PrintError(fmt.Sprintf("Stream error occurred: %v", err))
			}
		}
	}()

	StreamWithHighlighting(ctx, outChan)

	fmt.Printf("\n\n%s%s-------------------%s\n\n", Fg(51), Bold(), ResetStr())
}

func HandleInline(ctx context.Context, errormsg string, config cli.Config) {
	var stream *goai.TextStream
	var err error

	spinner := StartSpinner()

	switch config.Provider {
	case "Gemini":
		stream, err = gemini.Stream(ctx, config.APIKey, errormsg)
	case "Claude":
		stream, err = claude.Stream(ctx, config.APIKey, errormsg)
	case "ChatGPT":
		stream, err = chatgpt.Stream(ctx, config.APIKey, errormsg)
	default:
		spinner.Stop()
		PrintError("Provider not supported for Inline Mode")
		return
	}

	spinner.Stop()

	if err != nil {
		PrintError(fmt.Sprintf("Error initializing AI stream: %v", err))
		return
	}

	fmt.Printf("\n\n%s%s--- AI Analysis ---%s\n", Fg(51), Bold(), ResetStr())
	StreamWithHighlighting(ctx, stream.TextStream())
	fmt.Printf("\n\n%s%s-------------------%s\n\n", Fg(51), Bold(), ResetStr())

	if err := stream.Err(); err != nil {
		PrintError(fmt.Sprintf("Stream error occurred: %v", err))
	}
}
