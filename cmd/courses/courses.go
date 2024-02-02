package courses

import (
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"github.com/spf13/cobra"
)

var CoursesCmd = &cobra.Command{
	Use:     "courses",
	Short:   "Canvas Courses",
	Long:    "Commands for interacting with courses in Canvas",
	Aliases: []string{"course", "c"},
	Run:     execCoursesCmd,
}

func execCoursesCmd(cmd *cobra.Command, args []string) {
	log := util.Logger(cmd)
	client := util.Client(cmd)

	log.Debug().Str("canvasClient", client.String()).Msg("Contexts are cool")
}

func init() {
	// Flags, etc..
}
