package utils

import (
	"context"
	"errpipe/internal/ai/gemini"
	"errpipe/internal/cli"
	"fmt"

	"github.com/zendev-sh/goai"
)
var GEMINI_MODELS []string

func SendToAI(ctx context.Context, errormsg string, config cli.Config) {
	// Validate API Key for Inline CLI Mode, unless it's Free mode
	if config.APIKey == "" && config.Provider != "Free" {
		PrintError("API Key is required. Please run 'errpipe --init' to configure it.")
		return
	}

	switch config.Provider {
	case "Gemini":
		HandleInline(ctx, errormsg, config)
	default:
		PrintError(fmt.Sprintf("Provider %s is not supported.", config.Provider))
	}
}



func HandleInline(ctx context.Context, errormsg string, config cli.Config) {
	var stream *goai.TextStream
	var err error

	spinner := StartSpinner()

	switch config.Provider {
	case "Gemini":
		
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
