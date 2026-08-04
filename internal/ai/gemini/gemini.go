package gemini

import (
	"context"

	"github.com/zendev-sh/goai"
	"github.com/zendev-sh/goai/provider/google"
)

var GENERAL_INSTRUCTION = "Respond as short as possible with direct actionable fixes and a brief reason why. No fluff."
var FORMATTING_INSTRUCTION = "Format your response using markdown for readability:\n" +
	"- Wrap ALL code in fenced code blocks with a language tag (e.g. ```go).\n" +
	"- Use **bold** to highlight important terms, flags, or key concepts.\n" +
	"- Use `inline code` for short identifiers, filenames, or values inline with text.\n" +
	"- Use bullet lists for steps or multiple points.\n" +
	"- Do NOT use headers (#). Keep prose concise."
var SYSTEM_PROMPT = "Instruction: " + GENERAL_INSTRUCTION + "\nFormatting: " + FORMATTING_INSTRUCTION

func Stream(ctx context.Context, apikey, MODEL_ID string, errorMessage string) (*goai.TextStream, error){
	model := google.Chat(MODEL_ID, google.WithAPIKey(apikey))

	return goai.StreamText(ctx, model,
		goai.WithSystem(SYSTEM_PROMPT),
		goai.WithPrompt(errorMessage),
	)
}