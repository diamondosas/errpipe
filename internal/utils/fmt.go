package utils

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
)

// StreamWithHighlighting safely buffers text to look for markdown code blocks,
// buffers the inner code, and applies syntax highlighting using Chroma.
func StreamWithHighlighting(stream <-chan string) {
	inCodeBlock := false
	var textBuffer string
	var codeBuffer string
	var lang string

	for text := range stream {
		textBuffer += text

		for {
			if !inCodeBlock {
				// Look for start of code block
				idx := strings.Index(textBuffer, "```")
				if idx == -1 {
					// Safe to print everything EXCEPT the last 2 characters (in case they are part of an incoming "```")
					safeLen := len(textBuffer) - 2
					if safeLen > 0 {
						fmt.Print(Fg(255) + textBuffer[:safeLen] + ResetStr())
						textBuffer = textBuffer[safeLen:]
					}
					break // Break to wait for more chunks
				}

				// Print everything before ```
				if idx > 0 {
					fmt.Print(Fg(255) + textBuffer[:idx] + ResetStr())
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

	// Flush any remaining text/code when stream closes
	if inCodeBlock {
		codeBuffer += textBuffer
		highlightAndPrintCode(codeBuffer, lang)
	} else {
		if len(textBuffer) > 0 {
			fmt.Print(Fg(255) + textBuffer + ResetStr())
		}
	}
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
