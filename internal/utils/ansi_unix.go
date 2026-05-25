//go:build !windows

package utils

// EnableANSI enables virtual terminal processing on Windows to allow ANSI escape codes.
// On other platforms, it returns silently as ANSI is supported by default.
func EnableANSI() {
	// Not needed on non-Windows platforms
}
