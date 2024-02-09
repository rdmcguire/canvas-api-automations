package grades

import (
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"github.com/spf13/cobra"
)

var GradesCmd = &cobra.Command{
	Use:              "grades",
	Aliases:          []string{"g"},
	Short:            "sub-commands for grading",
	PersistentPreRun: util.EnsureCourseId,
}

func init() {
	GradesCmd.AddCommand(gradesDumpCmd)
	GradesCmd.AddCommand(gradesPushCmd)
	GradesCmd.AddCommand(gradesDumpCmd)

	GradesCmd.PersistentFlags().StringArray("email", []string{}, "Restrict to a provided email addresses")
}
