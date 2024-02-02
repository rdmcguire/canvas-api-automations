package cmd

import (
	"context"
	"os"
	"os/signal"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/assignments"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/courses"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/modules"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/students"
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
	level, _ := cmd.Flags().GetString("logLevel")
	logLevel := util.ParseLogLevel(level, defLogLevel)

	logger := log.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		Level(logLevel)

	cmd.SetContext(logger.WithContext(cmd.Context()))

	// Set globals
	log.Logger = logger
	zerolog.SetGlobalLevel(logLevel)

	util.SetClient(cmd, args)
}

func init() {
	// Add root-level flags
	rootCmd.PersistentFlags().StringP("logLevel", "l", "info",
		"Sets log level (fatal|error|warn|info|debug|trace)")

	// Register autocompletion funcs
	rootCmd.RegisterFlagCompletionFunc("logLevel", validLogLevels)

	// Add first-level sub-commands
	rootCmd.AddCommand(courses.CoursesCmd)
	rootCmd.AddCommand(assignments.AssignmentsCmd)
	rootCmd.AddCommand(students.StudentsCmd)
	rootCmd.AddCommand(modules.ModulesCmd)
}
