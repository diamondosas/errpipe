package claude

import (
	"errpipe/internal/utils/sys"
)

func OpenWeb(errorMessage string) {
	sys.OpenBrowser("Claude", errorMessage)
}
