package courses

import (
	"fmt"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"github.com/spf13/cobra"
)

var coursesShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"ls", "s"},
	Short:   "Show courses",
	Run:     execCoursesShowCmd,
}

func execCoursesShowCmd(cmd *cobra.Command, args []string) {
	log := util.Logger(cmd)
	client := util.Client(cmd)

	log.Trace().Msg("Random trace msg")

	fmt.Println("Courses:")
	for _, course := range client.ListCourses() {
		fmt.Println(canvas.CourseString(course))
	}
}

func init() {
	CoursesCmd.AddCommand(coursesShowCmd)
}
