package server

import (
	"github.com/emvi/logbuch"
)

// ConfigureLogging configures the logger.
func ConfigureLogging() {
	logbuch.SetFormatter(logbuch.NewFieldFormatter("2006-01-02T15:04:05", "\t"))

	if cfg.LogLevel == "debug" {
		logbuch.SetLevel(logbuch.LevelDebug)
	} else if cfg.LogLevel == "info" {
		logbuch.SetLevel(logbuch.LevelInfo)
	} else {
		logbuch.SetLevel(logbuch.LevelWarning)
	}
}
