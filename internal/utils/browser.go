package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
)

// findBrowser detects the default browser path per OS.
func findBrowser() (string, error) {
	switch runtime.GOOS {
	case "windows":
		// Read from registry
		cmd := exec.Command("reg", "query",
			`HKEY_CURRENT_USER\Software\Microsoft\Windows\Shell\Associations\UrlAssociations\http\UserChoice`,
			"/v", "ProgId",
		)
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}

		// Extract ProgId value (e.g., ChromeHTML, FirefoxURL)
		progID := extractRegValue(string(out))

		// Now find the actual executable for that ProgId
		cmd2 := exec.Command("reg", "query",
			fmt.Sprintf(`HKEY_CLASSES_ROOT\%s\shell\open\command`, progID),
		)
		out2, err := cmd2.Output()
		if err != nil {
			return "", err
		}
		return extractExePath(string(out2)), nil

	case "darwin":
		// macOS: use 'open' command or find via Launch Services
		// More reliable: use Python to query Launch Services
		cmd := exec.Command("python3", "-c", `
import subprocess
result = subprocess.run(
    ['osascript', '-e', 'POSIX path of (path to application id "com.apple.Safari")'],
    capture_output=True, text=True
)
print(result.stdout.strip())
`)
		out, err := cmd.Output()
		if err != nil {
			// Fallback to common paths
			return findBrowserInPaths([]string{
				"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
				"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
				"/Applications/Chromium.app/Contents/MacOS/Chromium",
				"/Applications/Firefox.app/Contents/MacOS/firefox",
				"/Applications/Safari.app/Contents/MacOS/Safari",
			})
		}
		return strings.TrimSpace(string(out)), nil

	case "linux":
		// Try xdg-settings first (most reliable)
		cmd := exec.Command("xdg-settings", "get", "default-web-browser")
		out, err := cmd.Output()
		if err == nil {
			desktopFile := strings.TrimSpace(string(out))
			// Resolve .desktop file to binary
			if path, err := resolveDesktopFile(desktopFile); err == nil {
				return path, nil
			}
		}

		// Fallback: check update-alternatives
		cmd2 := exec.Command("update-alternatives", "--query", "x-www-browser")
		out2, err := cmd2.Output()
		if err == nil {
			if path := extractAlternativePath(string(out2)); path != "" {
				return path, nil
			}
		}

		// Last resort: check common paths
		return findBrowserInPaths([]string{
			"/usr/bin/google-chrome",
			"/usr/bin/chromium-browser",
			"/usr/bin/chromium",
			"/usr/bin/firefox",
			"/snap/bin/chromium",
			"/usr/bin/microsoft-edge",
		})

	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// Helper: find first existing path from a list
func findBrowserInPaths(paths []string) (string, error) {
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("no browser found in common paths")
}

// Helper: extract value from Windows registry output
func extractRegValue(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "ProgId") || strings.Contains(line, "REG_SZ") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
		}
	}
	return ""
}

// Helper: extract .exe path from registry command output
func extractExePath(output string) string {
	start := strings.Index(output, `"`)
	end := strings.LastIndex(output, `"`)
	if start != -1 && end > start {
		return output[start+1 : end]
	}
	return ""
}

// Helper: resolve Linux .desktop file to binary path
func resolveDesktopFile(desktop string) (string, error) {
	searchPaths := []string{
		"/usr/share/applications/",
		"/usr/local/share/applications/",
		os.Getenv("HOME") + "/.local/share/applications/",
	}
	for _, dir := range searchPaths {
		path := dir + desktop
		if data, err := os.ReadFile(path); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(line, "Exec=") {
					exec := strings.TrimPrefix(line, "Exec=")
					exec = strings.Fields(exec)[0] // remove args like %u
					return exec, nil
				}
			}
		}
	}
	return "", fmt.Errorf("could not resolve .desktop file: %s", desktop)
}

func extractAlternativePath(output string) string {
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "Value:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Value:"))
		}
	}
	return ""
}


func OpenBrowser(provider, errorMessage string) {
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

	var pageURL, chatBoxSelector string

	switch provider {
	case "Gemini":
		pageURL = "https://gemini.google.com/app"
		chatBoxSelector = "rich-textarea"
	case "ChatGPT":
		pageURL = "https://chatgpt.com/"
		chatBoxSelector = "#prompt-textarea"
	case "Claude":
		pageURL = "https://claude.ai/new"
		chatBoxSelector = "[contenteditable='true']"
	default:
		pageURL = "https://gemini.google.com/app"
		chatBoxSelector = "rich-textarea"
	}

	page := browser.MustPage(pageURL)
	fmt.Printf("Waiting for %s chat box...\n", provider)

	chatBox := page.MustElement(chatBoxSelector)
	prompt := fmt.Sprintf("I'm getting the following error, please help me fix it:\n\n%s", errorMessage)
	chatBox.MustInput(prompt)
	page.Keyboard.Press(input.Enter)

	time.Sleep(10 * time.Second)

}
