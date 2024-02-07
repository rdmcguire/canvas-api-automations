package grades

import (
	"github.com/spf13/cobra"
)

var GradesCmd = &cobra.Command{
	Use:     "grades",
	Aliases: []string{"g"},
	Short:   "sub-commands for grading",
}

func init() {
	GradesCmd.AddCommand(gradesLoadCmd)
}
