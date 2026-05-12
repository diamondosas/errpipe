package gemini

import (
	"context"

	"github.com/zendev-sh/goai"
	"github.com/zendev-sh/goai/provider/google"
)

// Stream starts a streaming session with Gemini using GoAI
func Stream(apiKey, errorMessage string) (*goai.TextStream, error) {
	ctx := context.Background()
	model := google.Chat("gemini-2.5-pro", google.WithAPIKey(apiKey))

	return goai.StreamText(ctx, model,
		goai.WithSystem("You are an expert developer assistant. Respond in the shortest way possible with direct actionable fixes to the issue. No fluff"),
		goai.WithPrompt(errorMessage),
	)
}
