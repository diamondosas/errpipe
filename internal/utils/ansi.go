package utils

import (
	"fmt"
	"time"
)

// ── ANSI escape helpers ──────────────────────────────────────────────────────

func Fg(code int) string { return fmt.Sprintf("\033[38;5;%dm", code) }
func Bold() string       { return "\033[1m" }
func Dim() string        { return "\033[2m" }
func ResetStr() string   { return "\033[0m" }

// ── UI Components ────────────────────────────────────────────────────────────

type Spinner struct {
	stop chan struct{}
}

// StartSpinner initiates an animated spinner with the given message.
func StartSpinner(msg string) *Spinner {
	s := &Spinner{stop: make(chan struct{})}
	frames := []string{"⠋", "⠙", "⠸", "⠴", "⠦", "⠇"}

	go func() {
		i := 0
		for {
			select {
			case <-s.stop:
				// Clear the line when stopped
				fmt.Print("\r\033[K")
				return
			default:
				fmt.Printf("\r\t%s%s%s %s%s", Fg(214), frames[i%len(frames)], ResetStr(), Dim(), msg+ResetStr())
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return s
}

// Stop stops the spinner animation.
func (s *Spinner) Stop() {
	close(s.stop)
}

// PrintWelcome displays the initial banner with ANSI colors and formatting.
func PrintWelcome(provider, mode string) {
	banner := []string{
		`███████╗██████╗ ██████╗ ██████╗ ██╗██████╗ ███████╗`,
		`██╔════╝██╔══██╗██╔══██╗██╔══██╗██║██╔══██╗██╔════╝`,
		`█████╗  ██████╔╝██████╔╝██████╔╝██║██████╔╝█████╗  `,
		`██╔══╝  ██╔══██╗██╔══██╗██╔═══╝ ██║██╔═══╝ ██╔══╝  `,
		`███████╗██║  ██║██║  ██║██║     ██║██║     ███████╗`,
		`╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚═╝╚═╝     ╚══════╝`,
	}

	fmt.Println()
	for _, line := range banner {
		fmt.Println("\t" + Fg(51) + Bold() + line + ResetStr())
	}
	fmt.Println()
	fmt.Println("\t" + Fg(82) + Bold() + "» ERROR PIPE STARTED" + ResetStr())
	fmt.Println("\t" + Dim() + "Type 'errpipe --init' to setup application" + ResetStr())
	fmt.Printf("\t"+Fg(226)+"Using: %s (%s)\n\n"+ResetStr(), provider, mode)
}

// PrintPrompt prints the directory path prompt with ANSI styling.
func PrintPrompt(dir string) {
	fmt.Printf("%s[EP]%s %s%s%s> ", Fg(214)+Bold(), ResetStr(), Fg(39), dir, ResetStr())
}

// PrintAction prints a formatted, tab-indented action status (like sending to AI).
func PrintAction(action, detail string) {
	fmt.Printf("\n\t%s» %s:%s %s\n", Fg(208), action, ResetStr(), Dim()+detail+ResetStr())
}

// PrintSuccess prints a formatted success message
func PrintSuccess(msg string) {
	fmt.Printf("\t%s✔ %s%s\n\n", Fg(82), msg, ResetStr())
}

// PrintError prints a formatted error message
func PrintError(msg string) {
	fmt.Printf("\t%s✘ %s%s\n\n", Fg(196), msg, ResetStr())
}
