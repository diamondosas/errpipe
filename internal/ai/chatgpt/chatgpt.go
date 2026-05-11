package chatgpt

import (
	"context"

	"github.com/zendev-sh/goai"
	"github.com/zendev-sh/goai/provider/openai"
)

// Stream starts a streaming session with ChatGPT using GoAI
func Stream(apiKey, errorMessage string) (*goai.TextStream, error) {
	ctx := context.Background()
	model := openai.Chat("gpt-5.3-codex", openai.WithAPIKey(apiKey))

	return goai.StreamText(ctx, model,
		goai.WithSystem("You are an expert developer assistant. Analyze the error and provide a fix."),
		goai.WithPrompt(errorMessage),
	)
}
