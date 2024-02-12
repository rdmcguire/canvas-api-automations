package cmd

import (
	"context"
	"os"
	"os/signal"
	"strings"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/assignments"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/courses"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/grades"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/modules"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/students"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const defLogLevel = zerolog.InfoLevel

var rootCmd = &cobra.Command{
	Use:              "canvas-api-automations",
	PersistentPreRun: rootCmdPreRun,
	TraverseChildren: true,
	Short:            "Canvas API Interactions",
	Long: `Utilities for interacting with the canvas API
and brutal sorta-automations for Netacad courses`,
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
	commands := strings.Split(cmd.CommandPath(), " ")
	if len(commands) > 1 && commands[1] == "completion" {
		return
	}

	// Prepare logging
	level, _ := cmd.Flags().GetString("logLevel")
	logLevel := util.ParseLogLevel(level, defLogLevel)

	logger := log.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		Level(logLevel)
	ctx := logger.WithContext(cmd.Context())

	// Add any extra settings to context
	readOnly, _ := cmd.Flags().GetBool("readOnly")
	ctx = context.WithValue(ctx, "readOnly", readOnly)

	cmd.SetContext(ctx)

	// Set globals
	log.Logger = logger
	zerolog.SetGlobalLevel(logLevel)

	util.SetClient(cmd, args)
}

func init() {
	// Global settings
	cobra.EnableTraverseRunHooks = true // Allow all pre-run hooks to execute

	// Add root-level flags
	rootCmd.PersistentFlags().StringP("logLevel", "l", "info",
		"Sets log level (fatal|error|warn|info|debug|trace)")
	rootCmd.PersistentFlags().Int("courseID", 0,
		"Specify course ID, necessary for most sub-commands")
	rootCmd.PersistentFlags().Bool("readOnly", false,
		"Set to disable all non-GET http requests such as POST and PUT")

	// Register autocompletion funcs
	rootCmd.RegisterFlagCompletionFunc("logLevel", validLogLevels)
	rootCmd.RegisterFlagCompletionFunc("courseID", courses.ValidateCourseIdArg)

	// Add first-level sub-commands
	rootCmd.AddCommand(courses.CoursesCmd)
	rootCmd.AddCommand(assignments.AssignmentsCmd)
	rootCmd.AddCommand(students.StudentsCmd)
	rootCmd.AddCommand(modules.ModulesCmd)
	rootCmd.AddCommand(grades.GradesCmd)
	rootCmd.AddCommand(docsCmd)
}
