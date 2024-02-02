package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/courses"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const defLogLevel = zerolog.InfoLevel

var rootCmd = &cobra.Command{
	Use:   "canvas-api-automations",
	Short: "Canvas API Interactions",
	Long: `Utilities for interacting with the canvas API
and brutal sorta-automations for Netacad courses`,
	// Stores *canvas.Client in our global context
	PersistentPreRun: rootCmdPreRun,
}

func Execute() {
	ctx, cncl := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cncl()

	// Run cobra with our prepared context
	err := rootCmd.ExecuteContext(ctx)

	if err != nil {
		log.Error().Err(err)
		os.Exit(1)
	}
}

func rootCmdPreRun(cmd *cobra.Command, args []string) {
	// Prepare logging
	cmd.DebugFlags()
	level, _ := cmd.PersistentFlags().GetString("logLevel")
	logLevel := util.ParseLogLevel(level, defLogLevel)
	fmt.Println(level)
	fmt.Println(logLevel)

	logger := log.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		Level(logLevel)

	cmd.SetContext(logger.WithContext(cmd.Context()))

	// Set globals
	log.Logger = logger
	zerolog.SetGlobalLevel(logLevel)

	logger.Trace().Msg("Trace Logging")

	util.SetClient(cmd, args)
}

func init() {
	rootCmd.PersistentFlags().StringP("logLevel", "l", "info",
		"Sets log level (fatal|error|warn|info|debug|trace)")

	rootCmd.RegisterFlagCompletionFunc("logLevel", validLogLevels)

	// Add sub-commands
	rootCmd.AddCommand(courses.CoursesCmd)
}

func validLogLevels(cmd *cobra.Command, args []string, toComplete string,
) ([]string, cobra.ShellCompDirective) {
	return []string{"fatal", "error", "warn", "info", "debug", "trace"},
		cobra.ShellCompDirectiveNoFileComp
}
