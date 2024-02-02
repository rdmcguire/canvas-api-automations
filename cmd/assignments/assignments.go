package assignments

import (
	"github.com/spf13/cobra"
)

var AssignmentsCmd = &cobra.Command{
	Use:     "assignments",
	Short:   "Canvas Assignments",
	Long:    "Commands for interacting with assignments in Canvas",
	Aliases: []string{"assignment", "a"},
}

func init() {
	AssignmentsCmd.AddCommand(assignmentsShowCmd)
	AssignmentsCmd.AddCommand(assignmentsUpdateCmd)
}
