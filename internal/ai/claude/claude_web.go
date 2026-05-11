package claude

import (
	"errpipe/internal/utils"
)

func OpenWeb(errorMessage string) {
	utils.OpenBrowser("Claude", errorMessage)
}
