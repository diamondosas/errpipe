package cli

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

// Config represents the application configuration
type Config struct {
	Provider string
}

// InitApp starts the interactive setup process
func InitApp() {
	config := Config{}

	questions := []*survey.Question{
		{
			Name: "provider",
			Prompt: &survey.Select{
				Message: "Choose an AI provider:",
				Options: []string{"Claude", "ChatGPT", "Gemini"},
			},
		},
	}

	err := survey.Ask(questions, &config)
	if err != nil {
		fmt.Println("Setup cancelled.")
		return
	}

	fmt.Printf("\n✓ Selected provider: %s\n", config.Provider)
	
	// Start the process based on the selected provider
	startProcess(config.Provider)
}

func startProcess(provider string) {
	fmt.Printf("Starting process for %s...\n", provider)
	// Additional initialization logic for the provider will go here
}
