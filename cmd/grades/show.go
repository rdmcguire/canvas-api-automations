package grades

import "github.com/spf13/cobra"

var showGradesCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"list", "ls", "s", "dump"},
	Short:   "Shows graded submissions given any provided flags",
	Long:    "Run the show command with --email, --moduleID, or --assignmentID",
	Run:     execShowGradesCmd,
}

func execShowGradesCmd(cmd *cobra.Command, args []string) {
}
