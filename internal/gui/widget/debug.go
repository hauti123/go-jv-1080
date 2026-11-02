package widget

import (
	"log"

	"go-jv-1080/internal/util"
)

// debugEnabled turns on widget debug logging (drag handling, relayout,
// ...) when the LOG_GUI_WIDGETS environment variable is set.
var debugEnabled = util.EnvFlag("LOG_GUI_WIDGETS", false)

// debugf logs a debug message if LOG_GUI_WIDGETS is set.
func debugf(format string, args ...any) {
	if debugEnabled {
		log.Printf("WIDGET "+format, args...)
	}
}
