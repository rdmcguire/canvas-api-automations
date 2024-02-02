package assignments

import (
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/courses"
	"github.com/spf13/cobra"
)

func ValidateAssignmentUpdateArgs(cmd *cobra.Command, args []string, toComplete string,
) ([]string, cobra.ShellCompDirective) {
	if len(args) == 1 {
		return []string{}, cobra.ShellCompDirectiveFilterDirs
	} else if len(args) == 2 {
		return courses.ValidateCourseIdArg(cmd, args, toComplete)
	}
	return []string{}, cobra.ShellCompDirectiveNoFileComp
}
