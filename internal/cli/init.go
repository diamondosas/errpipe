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
		fmt.Printf("\n  [33m✔ Tip:[0m Copy your API key first, then [1mright-click[0m inside the terminal to paste it.\n\n")
		apiKeyQuestion := &survey.Input{
			Message: fmt.Sprintf("Enter your %s API Key:", config.Provider),
			Help:    "Right-click inside the terminal window to paste your API key. On Linux/Mac use Shift+Insert or Ctrl+Shift+V.",
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
