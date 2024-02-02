package util

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// Retrieves the logger from the command context
func Logger(cmd *cobra.Command) *zerolog.Logger {
	return zerolog.Ctx(cmd.Context())
}

func ParseLogLevel(level string, defaultLevel zerolog.Level) zerolog.Level {
	switch strings.ToLower(level) {
	case "fatal":
		return zerolog.FatalLevel
	case "error":
		return zerolog.ErrorLevel
	case "warn":
		return zerolog.WarnLevel
	case "info":
		return zerolog.InfoLevel
	case "debug":
		return zerolog.DebugLevel
	case "trace":
		return zerolog.TraceLevel
	}
	return defaultLevel
}
