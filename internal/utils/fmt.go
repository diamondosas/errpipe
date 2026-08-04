package utils

import (
	"context"
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
)

// StreamWithHighlighting safely buffers text to look for markdown code blocks,
// buffers the inner code, and applies syntax highlighting using Chroma.
func StreamWithHighlighting(ctx context.Context, stream <-chan string) {
	inCodeBlock := false
	var textBuffer string
	var codeBuffer string
	var lang string

	for {
		select {
		case <-ctx.Done():
			fmt.Print(ResetStr()) // Ensure we don't leave bad formatting
			fmt.Printf("\n\n%s[!] Analysis cancelled.%s\n", Fg(196), ResetStr())
			return
		case text, ok := <-stream:
			if !ok {
				// Flush any remaining text/code when stream closes
				if inCodeBlock {
					codeBuffer += textBuffer
					highlightAndPrintCode(codeBuffer, lang)
				} else {
					if len(textBuffer) > 0 {
						fmt.Print(formatInlineMarkdown(textBuffer))
					}
				}
				return
			}
			textBuffer += text

			for {
				if !inCodeBlock {
					// Look for start of code block
					idx := strings.Index(textBuffer, "```")
					if idx == -1 {
						// Safe to print everything EXCEPT the last 2 characters (in case they are part of an incoming "```")
						safeLen := len(textBuffer) - 2
						if safeLen > 0 {
							fmt.Print(formatInlineMarkdown(textBuffer[:safeLen]))
							textBuffer = textBuffer[safeLen:]
						}
						break // Break to wait for more chunks
					}

					// Print everything before ```
					if idx > 0 {
						fmt.Print(formatInlineMarkdown(textBuffer[:idx]))
					}
					textBuffer = textBuffer[idx+3:]

					// Look for the end of the line to extract the language
					nlIdx := strings.Index(textBuffer, "\n")
					if nlIdx == -1 {
						// We haven't received the newline after ``` yet.
						// Put "```" back into the buffer and wait for more chunks.
						textBuffer = "```" + textBuffer
						break
					}

					// Extract language and skip the newline
					lang = strings.TrimSpace(textBuffer[:nlIdx])
					textBuffer = textBuffer[nlIdx+1:]

					inCodeBlock = true
					codeBuffer = ""
				} else {
					// We are in a code block, look for closing ```
					idx := strings.Index(textBuffer, "```")
					if idx == -1 {
						// Safe to buffer everything EXCEPT the last 2 characters (in case they are part of an incoming "```")
						safeLen := len(textBuffer) - 2
						if safeLen > 0 {
							codeBuffer += textBuffer[:safeLen]
							textBuffer = textBuffer[safeLen:]
						}
						break // Break to wait for more chunks
					}

					// Found closing ```
					codeBuffer += textBuffer[:idx]

					// Highlight and print the buffered code!
					highlightAndPrintCode(codeBuffer, lang)

					// Reset state
					inCodeBlock = false
					codeBuffer = ""
					textBuffer = textBuffer[idx+3:]
				}
			}
		}
	}
}

// formatInlineMarkdown applies ANSI styling for inline markdown produced by the AI:
//   - ***text***  → bold + italic
//   - **text**    → bold
//   - *text*      → italic
//   - `text`      → cyan  (inline code)
//   - - item / * item → styled bullet point
//
// Patterns are applied in specificity order (most specific first) to avoid
// partial matches (e.g. *** being confused with **).
var (
	reBoldItalic = regexp.MustCompile(`\*{3}([^*]+)\*{3}`)
	reBold       = regexp.MustCompile(`\*{2}([^*]+)\*{2}`)
	reItalic     = regexp.MustCompile(`\*([^*\n]+)\*`)
	reInlineCode = regexp.MustCompile("`([^`\n]+)`")
)

// ansi escape helpers
const (
	ansiBold      = "\033[1m"
	ansiItalic    = "\033[3m"
	ansiCyan      = "\033[36m"
	ansiDim       = "\033[2m"
	ansiReset     = "\033[0m"
)

func formatInlineMarkdown(s string) string {
	// Process line-by-line so bullet detection is scoped correctly
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Detect bullet list items: lines starting with "- " or "* " (but not "**")
		isBullet := (strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ")) &&
			!strings.HasPrefix(trimmed, "**")
		if isBullet {
			// Replace the leading "- " / "* " with a styled bullet glyph
			rest := trimmed[2:]
			line = ansiReset + Fg(240) + "  •" + ansiReset + " " + rest
		}

		// Apply inline patterns in specificity order
		line = reBoldItalic.ReplaceAllStringFunc(line, func(m string) string {
			inner := reBoldItalic.FindStringSubmatch(m)[1]
			return ansiReset + ansiBold + ansiItalic + inner + ansiReset
		})
		line = reBold.ReplaceAllStringFunc(line, func(m string) string {
			inner := reBold.FindStringSubmatch(m)[1]
			return ansiReset + ansiBold + inner + ansiReset
		})
		line = reItalic.ReplaceAllStringFunc(line, func(m string) string {
			inner := reItalic.FindStringSubmatch(m)[1]
			return ansiReset + ansiItalic + inner + ansiReset
		})
		line = reInlineCode.ReplaceAllStringFunc(line, func(m string) string {
			inner := reInlineCode.FindStringSubmatch(m)[1]
			return ansiReset + ansiCyan + inner + ansiReset
		})

		lines[i] = line
	}

	return Fg(255) + strings.Join(lines, "\n") + ResetStr()
}

// highlightAndPrintCode formats the buffered code block and syntax highlights it.
func highlightAndPrintCode(code, lang string) {
	if lang == "" {
		lang = "text" // Fallback language
	}

	var buf bytes.Buffer
	// Use Chroma to highlight ("terminal256" outputs ANSI color codes, "monokai" is a nice theme)
	err := quick.Highlight(&buf, code, lang, "terminal256", "monokai")
	if err != nil {
		// Fallback to plain ANSI cyan if Chroma fails
		fmt.Printf("\n\033[36m%s\033[0m\n", code)
		return
	}

	// Print with an indentation and vertical bar to visually separate it from normal text
	fmt.Println()
	highlighted := buf.String()
	lines := strings.Split(strings.TrimSuffix(highlighted, "\n"), "\n")
	for _, line := range lines {
		// Reset string `\033[0m` is added to ensure colors don't bleed into the margin
		fmt.Printf("    \033[90m|\033[0m %s\n", line)
	}
	fmt.Println()
}
