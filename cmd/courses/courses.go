package courses

import (
	"github.com/spf13/cobra"
)

var CoursesCmd = &cobra.Command{
	Use:     "courses",
	Short:   "Canvas Courses",
	Long:    "Commands for interacting with courses in Canvas",
	Aliases: []string{"course", "c"},
}

func init() {
	CoursesCmd.AddCommand(coursesShowCmd)
}
