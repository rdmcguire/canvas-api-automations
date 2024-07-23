package assignments

import (
	"strings"

	"github.com/spf13/cobra"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/courses"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
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

func ValidAssignmentIDArg(cmd *cobra.Command, args []string, toComplete string,
) ([]string, cobra.ShellCompDirective) {
	completions := make([]string, 0)
	assignments, err := util.Client(cmd).ListAssignments(util.GetCourseIdStr(cmd))
	if err != nil || len(assignments) < 1 {
		util.Logger(cmd).Error().Err(err).Msg("Failed to list assignments")
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	for _, a := range assignments {
		id := canvas.StrOrNil(a.Id)
		if toComplete == "" || strings.HasPrefix(id, toComplete) {
			completions = append(completions, id)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}
