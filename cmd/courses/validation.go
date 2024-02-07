package courses

import (
	"strings"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"github.com/spf13/cobra"
)

func ValidateCourseIdArg(cmd *cobra.Command, args []string, toComplete string,
) ([]string, cobra.ShellCompDirective) {
	// First retrieve all courses
	client := util.Client(cmd)
	courses := client.ListCourses()

	// Then filter and return
	validCourses := make([]string, 0, len(courses))
	for _, c := range courses {
		if strings.HasPrefix(canvas.StrOrNil(c.Id), toComplete) {
			validCourses = append(validCourses, canvas.StrOrNil(c.Id))
		}
	}

	return validCourses, cobra.ShellCompDirectiveNoFileComp
}
