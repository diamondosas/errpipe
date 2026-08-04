package cli

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

// InitApp starts the interactive setup process
func InitApp() {
	// Try to load existing config
	existingConfig, err := LoadConfig()
	if err == nil {
		fmt.Printf("Current Configuration: AI - %s\n\n", existingConfig.Provider)
	}

	config := Config{}

	// First question: Choose Provider
	providerQuestion := &survey.Select{
		Message: "Choose an AI provider:",
		Options: []string{"Gemini", "Ollama"},
	}

	err = survey.AskOne(providerQuestion, &config.Provider)
	if err != nil {
		fmt.Println("Setup cancelled.")
		return
	}


	// Ask for API Key if provider is not Free
	if config.Provider != "Free" {
		apiKeyQuestion := &survey.Input{
			Message: fmt.Sprintf("Enter your %s API Key:", config.Provider),
		}
		err = survey.AskOne(apiKeyQuestion, &config.APIKey)
		if err != nil {
			fmt.Println("Setup cancelled.")
			return
		}
	}

	// Save the configuration
	err = SaveConfig(config)
	if err != nil {
		fmt.Printf("Error saving configuration: %v\n", err)
	} else {
		fmt.Printf("\n✓ Configuration Saved: %s\n", config.Provider)
	}
}
