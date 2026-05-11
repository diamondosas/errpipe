package chatgpt

import (
	"errpipe/internal/utils"
)

func OpenWeb(errorMessage string) {
	utils.OpenBrowser("ChatGPT", errorMessage)
}
