package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// PantryData represents the JSON structure from the pantry server
type Response struct {
	One   string `json:"1"`
	Two   string `json:"2"`
	Three string `json:"3"`
}

func init() {
	go FetchModels()
}

var ModelsPath string
var response Response

// FetchModels fetches the models from pantry or loads them from the local app data folder cache.
func FetchModels() {
	// Make HTTP request with timeout to fetch fresh models
	req, err := http.NewRequest("GET", "https://api.jsonbin.io/v3/b/6a71cf20da38895dfeb80b24?meta=false", nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", "$2a$10$NxHQ2PNwBeufrft4HlJTGOSJ4UTV.8T9da3RrE3cGbU3aVwtwR2gu")
	//Dont bother trying to Hack Because it is a Read Only Access Key
	req.Header.Set("X-Access-Key", "$2a$10$SYI3HJT7Xi1z6siLVLP2Pec0cHhNJIqvKPVNkZcWEMWp8J9EQ6HD6")

	cli := &http.Client{Timeout: 10 * time.Second}
	resp, err := cli.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return
	}

	models := []string{response.One, response.Two, response.Three}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return
	}
	errpipeDir := filepath.Join(configDir, "errpipe")
	ModelsPath = filepath.Join(errpipeDir, "models.json")
	// Ensure directory exists
	_ = os.MkdirAll(errpipeDir, 0755)

	if cachedData, err := json.Marshal(models); err == nil {
		_ = os.WriteFile(ModelsPath, cachedData, 0644)
	}
}

func GetModels() []string {
	// Try to load cached models first so they are immediately available
	if data, err := os.ReadFile(ModelsPath); err == nil {
		var cachedModels []string
		if err := json.Unmarshal(data, &cachedModels); err == nil && len(cachedModels) > 0 {
			return cachedModels
		}
	}
	return nil
}
