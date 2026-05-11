package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	Provider string `json:"provider"`
	Mode     string `json:"mode"`
	APIKey   string `json:"api_key"`
}

// GetConfigPath returns the path to the config.json file
func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	errpipeDir := filepath.Join(configDir, "errpipe")
	if _, err := os.Stat(errpipeDir); os.IsNotExist(err) {
		err = os.MkdirAll(errpipeDir, 0755)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(errpipeDir, "config.json"), nil
}

// SaveConfig saves the configuration to disk
func SaveConfig(config Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// LoadConfig loads the configuration from disk
func LoadConfig() (Config, error) {
	var config Config
	path, err := GetConfigPath()
	if err != nil {
		return config, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}
