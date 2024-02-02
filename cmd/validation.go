package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

func validLogLevels(cmd *cobra.Command, args []string, toComplete string,
) ([]string, cobra.ShellCompDirective) {
	levels := []string{"fatal", "error", "warn", "info", "debug", "trace"}
	logLevels := make([]string, 0)
	for _, level := range levels {
		if strings.HasPrefix(level, toComplete) {
			logLevels = append(logLevels, level)
		}
	}
	return logLevels, cobra.ShellCompDirectiveNoFileComp
}
