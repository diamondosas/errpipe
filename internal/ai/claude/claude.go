package claude

import (
	"context"

	"github.com/zendev-sh/goai"
	"github.com/zendev-sh/goai/provider/anthropic"
)

// Stream starts a streaming session with Claude using GoAI
func Stream(ctx context.Context, apiKey, errorMessage string) (*goai.TextStream, error) {
	model := anthropic.Chat("claude-3-5-haiku-latest", anthropic.WithAPIKey(apiKey))

	return goai.StreamText(ctx, model,
		goai.WithSystem("You are an expert developer assistant. Respond in the shortest way possible with direct actionable fixes to the issue. No fluff"),
		goai.WithPrompt(errorMessage),
	)
}
