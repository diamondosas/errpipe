package free

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// Models are the available models for DuckDuckGo chat
var Models = []string{""}

func setBaseHeaders(req *http.Request) {
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "fr-FR,fr;q=0.6")
	req.Header.Set("Cache-Control", "no-store")
	req.Header.Set("DNT", "1")
	req.Header.Set("Referer", "https://duckduckgo.com/")
	req.Header.Set("Sec-CH-UA", `"Not)A;Brand";v="8", "Chromium";v="138", "Brave";v="138"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
}

func getVQD(client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", "https://duckduckgo.com/duckchat/v1/status", nil)
	if err != nil {
		return "", err
	}
	setBaseHeaders(req)
	req.Header.Set("x-vqd-accept", "1")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code: %d", resp.StatusCode)
	}

	vqd := resp.Header.Get("x-vqd-4")
	if vqd == "" {
		vqd = resp.Header.Get("x-vqd-hash-1")
	}
	if vqd == "" {
		return "", fmt.Errorf("no VQD token found")
	}
	return vqd, nil
}

// StreamToChan streams the AI response chunks to outChan and any error to errChan.
func StreamToChan(errorMessage string, outChan chan<- string, errChan chan<- error) {
	defer close(outChan)
	defer close(errChan)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse("https://duckduckgo.com")
	jar.SetCookies(u, []*http.Cookie{
		{Name: "5", Value: "1"},
		{Name: "dcm", Value: "3"},
		{Name: "dcs", Value: "1"},
	})

	client := &http.Client{
		Jar: jar,
	}

	vqd, err := getVQD(client)
	if err != nil {
		errChan <- fmt.Errorf("failed to get VQD: %w", err)
		return
	}

	modelsToTry := Models
	if len(modelsToTry) == 0 {
		modelsToTry = []string{""}
	}

	var lastErr error
	success := false

	for _, model := range modelsToTry {
		payload := map[string]interface{}{
			"model":  model,
			"metadata": map[string]interface{}{
				"toolChoice": map[string]bool{
					"NewsSearch":      false,
					"VideosSearch":    false,
					"LocalSearch":     false,
					"WeatherForecast": false,
				},
			},
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": "You are an expert developer assistant. Respond in the shortest way possible with direct actionable fixes to the issue. No fluff. Issue: " + errorMessage,
				},
			},
			"canUseTools":          false,
			"canUseApproxLocation": false,
		}

		bodyData, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", "https://duckduckgo.com/duckchat/v1/chat", bytes.NewReader(bodyData))
		if err != nil {
			lastErr = err
			continue
		}

		setBaseHeaders(req)
		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://duckduckgo.com")
		req.Header.Set("x-vqd-4", vqd)
		req.Header.Set("x-vqd-hash-1", "eyJzZXJ2ZXJfaGFzaGVzIjpbImRQSlJJTWczZnFYQXIvaStaa3c2cEpFVzEwckdTdmxJVlVkNlFsOVRGWXc9IiwiMUN3Qzg3N0Q3WXE1dzlEeTc4UjhBVi9qZVZWaUlYbmV0Q0xvckx3c01QZz0iLCJQSzc3TGc2L25weDdWQ2J2UWxsTEhBR3cyenJIVmEvQUFBRFBhQTl1ekVRPSJdLCJjbGllbnRfaGFzaGVzIjpbImxWblI0MStCMVFWZ0o4d0hhMUdBNmdxR0JoSjlWdjN5K0dISkdGekJmTGM9IiwiVS9RRUc2RE1qdEU4V2hHU1FxOUU1Z0VGNmw1SWJrNk9NVlBuY01DU1licz0iLCJ6SURsYUNvZG9JUjNwbTNSVTlWOUJXaUJkZDJqenRMODAyN0VYTHhkWll3PSJdLCJzaWduYWxzIjp7fSwibWV0YSI6eyJ2IjoiNCIsImNoYWxsZW5nZV9pZCI6ImM4M2Q0ZTc5NTU2MjJmZjU3Mzc0ZDUzOTk2ZjliMmJhZGE2ZDQxZTMzNDM1ZjVlNzMyYjFmNmZjNmQ0ZTE1NzVoOGpidCIsInRpbWVzdGFtcCI6IjE3NTIxNTU3Nzc4NjYiLCJvcmlnaW4iOiJodHRwczovL2R1Y2tkdWNrZ28uY29tIiwic3RhY2siOiJFcnJvclxuYXQgRSAoaHR0cHM6Ly9kdWNrZHVja2dvLmNvbS9kaXN0L3dwbS5jaGF0LjcwZWFjYTZhZWEyOTQ4YjBiYjYwLmpzOjE6MTQ4MjUpXG5hdCBhc3luYyBodHRwczovL2R1Y2tkdWNrZ28uY29tIiwic3RhY2siOiJFcnJvclxuYXQgRSAoaHR0cHM6Ly9kdWNrZHVja2dvLmNvbS9kaXN0L3dwbS5jaGF0LjcwZWFjYTZhZWEyOTQ4YjBiYjYwLmpzOjE6MTQ4MjUpXG5hdCBhc3luYyBodHRwczovL2R1Y2tkdWNrZ28uY29tL2Rpc3Qvd3BtLmNoYXQuNzBlYWNhNmFlYTI5NDhiMGJiNjAuanM6MToxNjk4NSIsImR1cmF0aW9uIjoiNTgifX0=")
		req.Header.Set("x-fe-signals", "eyJzdGFydCI6MTc1MjE1NTc3NzQ4MCwiZXZlbnRzIjpbeyJuYW1lIjoic3RhcnROZXdDaGF0IiwiZGVsdGEiOjc1fSx7Im5hbWUiOiJyZWNlbnRDaGF0c0xpc3RJbXByZXNzaW9uIiwiZGVsdGEiOjEyNH1dLCJlbmQiOjQzNDN9")
		req.Header.Set("x-fe-version", "serp_20250710_090702_ET-70eaca6aea2948b0bb60")

		resp, err := client.Do(req)
		if err != nil {
			// This is a network error (e.g. DNS failure, connection timed out)
			// Do not retry next models, send error to errChan and return.
			errChan <- err
			return
		}

		if resp.StatusCode != 200 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("chat API returned %d for model %q: %s", resp.StatusCode, model, string(bodyBytes))
			continue
		}

		success = true
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				resp.Body.Close()
				if err != io.EOF {
					errChan <- err
				}
				return
			}

			lineStr := strings.TrimSpace(string(line))
			if !strings.HasPrefix(lineStr, "data: ") {
				continue
			}

			data := strings.TrimPrefix(lineStr, "data: ")
			if data == "[DONE]" || data == "" {
				resp.Body.Close()
				return
			}

			var chunk map[string]interface{}
			if err := json.Unmarshal([]byte(data), &chunk); err == nil {
				if msg, ok := chunk["message"].(string); ok && msg != "" {
					outChan <- msg
				}
			}
		}
	}

	if !success {
		if lastErr != nil {
			errChan <- lastErr
		} else {
			errChan <- fmt.Errorf("no models succeeded")
		}
	}
}
