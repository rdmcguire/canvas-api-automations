package students

import (
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"github.com/spf13/cobra"
)

var StudentsCmd = &cobra.Command{
	Use:              "students",
	Aliases:          []string{"s"},
	Short:            "Sub-command for student actions",
	PersistentPreRun: util.EnsureCourseId,
}

func init() {
	StudentsCmd.AddCommand(studentsExportCmd)
}
