package assignments

import (
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"github.com/spf13/cobra"
)

var AssignmentsCmd = &cobra.Command{
	Use:              "assignments",
	Short:            "Canvas Assignments",
	Long:             "Commands for interacting with assignments in Canvas",
	Aliases:          []string{"assignment", "a"},
	PersistentPreRun: util.EnsureCourseId,
}

func init() {
	AssignmentsCmd.AddCommand(assignmentsShowCmd)
	AssignmentsCmd.AddCommand(assignmentsUpdateCmd)
}
