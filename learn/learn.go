package main

import (
    "fmt"
    "github.com/AlecAivazis/survey/v2"
)

type ProjectConfig struct {
    Name       string
    Framework  string
    UseTS      bool
    PackageMan string
}

func main() {
    config := ProjectConfig{}

    // Each question is defined separately
    questions := []*survey.Question{
        {
            Name: "name",
            Prompt: &survey.Input{
                Message: "What is your project name?",
                Default: "my-app",
            },
            Validate: survey.Required, // built-in validator
        },
        {
            Name: "framework",
            Prompt: &survey.Select{
                Message: "Choose a framework:",
                Options: []string{"Next.js", "Remix", "Astro"},
            },
        },
        {
            Name: "useTS",
            Prompt: &survey.Confirm{
                Message: "Use TypeScript?",
                Default: true,
            },
        },
        {
            Name: "packageMan",
            Prompt: &survey.Select{
                Message: "Package manager:",
                Options: []string{"npm", "pnpm", "bun"},
            },
        },
    }

    // survey fills the struct for you via field name matching
    err := survey.Ask(questions, &config)
    if err != nil {
        fmt.Println("Cancelled.")
        return
    }

    fmt.Printf("\n✓ Creating %s with %s (%s)\n", config.Name, config.Framework, config.PackageMan)
    if config.UseTS {
        fmt.Println("✓ TypeScript enabled")
    }
}