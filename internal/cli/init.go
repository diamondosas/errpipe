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
		fmt.Printf("Current Configuration: AI - %s MODE - %s\n\n", existingConfig.Provider, existingConfig.Mode)
	}

	config := Config{}

	// First question: Choose Provider
	providerQuestion := &survey.Select{
		Message: "Choose an AI provider:",
		Options: []string{"Gemini", "Claude", "ChatGPT"},
	}

	err = survey.AskOne(providerQuestion, &config.Provider)
	if err != nil {
		fmt.Println("Setup cancelled.")
		return
	}

	// Second question: Choose Mode based on provider
	modeQuestion := &survey.Select{
		Message: fmt.Sprintf("Choose mode for %s:", config.Provider),
		Options: []string{"Inline CLI Mode", config.Provider + " CLI Mode", "Web Mode"},
	}

	err = survey.AskOne(modeQuestion, &config.Mode)
	if err != nil {
		fmt.Println("Setup cancelled.")
		return
	}

	// Ask for API Key if Inline CLI Mode is selected
	if config.Mode == "Inline CLI Mode" {
		apiKeyQuestion := &survey.Input{
			Message: "Enter your API Key:",
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
		fmt.Printf("\n✓ Configuration Saved: %s in %s\n", config.Provider, config.Mode)
	}

	// Start the process based on the selected provider and mode
	startProcess(config.Provider, config.Mode)
}

func startProcess(provider, mode string) {
	fmt.Printf("Initializing %s (%s)...\n", provider, mode)
	// Additional initialization logic for the provider and mode will go here
}
