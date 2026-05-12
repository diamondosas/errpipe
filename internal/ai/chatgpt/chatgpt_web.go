package chatgpt

import (
	"errpipe/internal/utils/sys"
)

func OpenWeb(errorMessage string) {
	sys.OpenBrowser("ChatGPT", errorMessage)
}
