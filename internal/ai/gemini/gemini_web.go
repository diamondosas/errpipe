package gemini

import (
	"errpipe/internal/utils/sys"
)

func OpenWeb(errorMessage string) {
	sys.OpenBrowser("Gemini", errorMessage)
}
