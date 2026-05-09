package gemini

import (
	"errpipe/internal/utils"
)

func OpenWeb(errorMessage string) {
	utils.OpenBrowser(errorMessage)
}
