package students

import "github.com/spf13/cobra"

var StudentsCmd = &cobra.Command{
	Use:     "students",
	Aliases: []string{"s"},
	Short:   "Sub-command for student actions",
}

func init() {
	StudentsCmd.AddCommand(studentsExportCmd)
}
