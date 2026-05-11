package claude

import (
	"context"

	"github.com/zendev-sh/goai"
	"github.com/zendev-sh/goai/provider/anthropic"
)

// Stream starts a streaming session with Claude using GoAI
func Stream(apiKey, errorMessage string) (*goai.TextStream, error) {
	ctx := context.Background()
	model := anthropic.Chat("claude-opus-4.6", anthropic.WithAPIKey(apiKey))

	return goai.StreamText(ctx, model,
		goai.WithSystem("You are an expert developer assistant. Analyze the error and provide a fix."),
		goai.WithPrompt(errorMessage),
	)
}
