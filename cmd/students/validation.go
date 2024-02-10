package students

import (
	"strings"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"github.com/spf13/cobra"
)

func ValidateEmailArg(cmd *cobra.Command, args []string, toComplete string,
) ([]string, cobra.ShellCompDirective) {
	// First retrieve all courses
	client := util.Client(cmd)
	users := client.ListUsersInCourse(util.GetCourseIdStr(cmd), "")

	// Then filter and return
	validUsers := make([]string, 0, len(users))
	for _, user := range users {
		if strings.HasPrefix(canvas.StrOrNil(user.Email), toComplete) {
			validUsers = append(validUsers, canvas.StrOrNil(user.Email))
		}
	}

	return validUsers, cobra.ShellCompDirectiveNoFileComp
}
