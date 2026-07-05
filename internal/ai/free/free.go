package free

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Models is defined to allow utils/models.go to compile, but is ignored.
var Models = []string{
	"gpt-4o-mini",
	"meta-llama/Llama-3.3-70B-Instruct-Turbo",
	"claude-3-haiku-20240307",
}

// This structure matches the JSON boxes the AI sends over the stream
type ChatChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

// StreamToChan streams the AI response chunks to outChan and any error to errChan.
func StreamToChan(errorMessage string, outChan chan<- string, errChan chan<- error) {
	defer close(outChan)
	defer close(errChan)

	prompt := "You are an expert developer assistant. Respond in the shortest way possible with direct actionable fixes to the issue. No fluff. Issue: " + errorMessage
	encodedPrompt := url.PathEscape(prompt)
	apiURL := "https://text.pollinations.ai/" + encodedPrompt + "?stream=true"

	response, err := http.Get(apiURL)
	if err != nil {
		errChan <- fmt.Errorf("could not reach AI: %w", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(response.Body)
		errChan <- fmt.Errorf("status code %d: %s", response.StatusCode, string(bodyBytes))
		return
	}

	reader := bufio.NewReader(response.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				errChan <- err
			}
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			jsonData := strings.TrimPrefix(line, "data: ")

			if jsonData == "[DONE]" {
				break
			}

			var chunk ChatChunk
			if err := json.Unmarshal([]byte(jsonData), &chunk); err == nil {
				if len(chunk.Choices) > 0 {
					content := chunk.Choices[0].Delta.Content
					if content != "" {
						outChan <- content
					}
				}wwww
			}
		}
	}
}
