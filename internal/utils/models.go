package utils

import (
	"encoding/json"
	"errpipe/internal/ai/free"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// PantryData represents the JSON structure from the pantry server
type PantryData struct {
	One   string `json:"1"`
	Two   string `json:"2"`
	Three string `json:"3"`
}

func init() {
	go FetchModels()
}

// FetchModels fetches the models from pantry or loads them from the local app data folder cache.
func FetchModels() {
	// Get the path to models.json in user's config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		return
	}
	errpipeDir := filepath.Join(configDir, "errpipe")
	modelsPath := filepath.Join(errpipeDir, "models.json")

	// Try to load cached models first so they are immediately available
	if data, err := os.ReadFile(modelsPath); err == nil {
		var cachedModels []string
		if err := json.Unmarshal(data, &cachedModels); err == nil && len(cachedModels) > 0 {
			free.Models = cachedModels
		}
	}

	// Make HTTP request with timeout to fetch fresh models
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://getpantry.cloud/apiv1/public/ec41d9bde014139aa3de25a72a090a89")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return
	}

	var pData PantryData
	if err := json.NewDecoder(resp.Body).Decode(&pData); err != nil {
		return
	}

	// Construct models list from pantry data
	var models []string
	if pData.One != "" {
		models = append(models, pData.One)
	}
	if pData.Two != "" {
		models = append(models, pData.Two)
	}
	if pData.Three != "" {
		models = append(models, pData.Three)
	}

	if len(models) == 0 {
		return
	}

	// Update free.Models and save to cache file
	free.Models = models

	// Ensure directory exists
	_ = os.MkdirAll(errpipeDir, 0755)

	if cachedData, err := json.Marshal(models); err == nil {
		_ = os.WriteFile(modelsPath, cachedData, 0644)
	}
}
