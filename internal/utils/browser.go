package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
)

// findBrowser detects the default Chromium-based browser path per OS.
func findBrowser() (string, error) {
	var candidates []string

	switch runtime.GOOS {
	case "windows":
		// Use environment variables instead of hardcoded paths
		programFiles := os.Getenv("ProgramFiles")
		programFilesX86 := os.Getenv("ProgramFiles(x86)")
		localAppData := os.Getenv("LOCALAPPDATA")

		candidates = []string{
			localAppData + `\Microsoft\Edge\Application\msedge.exe`,
			programFilesX86 + `\Microsoft\Edge\Application\msedge.exe`,
			programFiles + `\Microsoft\Edge\Application\msedge.exe`,
			localAppData + `\Google\Chrome\Application\chrome.exe`,
			programFiles + `\Google\Chrome\Application\chrome.exe`,
			programFilesX86 + `\Google\Chrome\Application\chrome.exe`,
		}

	case "darwin":
		candidates = []string{
			"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
		}

	case "linux":
		// On Linux these are on PATH, so use exec.LookPath
		for _, name := range []string{"microsoft-edge", "google-chrome", "chromium-browser", "chromium"} {
			if path, err := exec.LookPath(name); err == nil {
				return path, nil
			}
		}
		return "", fmt.Errorf("no supported browser found on Linux")

	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Check each candidate path and return the first one that exists
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("no supported browser found on %s", runtime.GOOS)
}


func openBroswer(errorMessage string){
	browserPath, err := findBrowser()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	wsURL, err := launcher.NewUserMode().Bin(browserPath).Launch()
	if err != nil {
		fmt.Printf("Failed to connect to browser: %v\n", err)
		os.Exit(1)
	}

	browser := rod.New().ControlURL(wsURL).MustConnect()
	page := browser.MustPage("https://gemini.google.com/app")
	fmt.Println("Waiting for chat box...")

	chatBox := page.MustElement("rich-textarea")
	prompt := fmt.Sprintf("I'm getting the following error, please help me fix it:\n\n%s", errorMessage)
	chatBox.MustInput(prompt)
	page.Keyboard.Press(input.Enter)

	time.Sleep(10 * time.Second)

}
