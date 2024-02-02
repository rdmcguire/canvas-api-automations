package modules

import (
	"fmt"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/courses"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"github.com/spf13/cobra"
)

var modulesShowCmd = &cobra.Command{
	Use:               "show",
	Aliases:           []string{"ls", "s"},
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: courses.ValidateCourseIdArg,
	Run:               execModulesShowCmd,
	Short:             "Show modules",
}

func execModulesShowCmd(cmd *cobra.Command, args []string) {
	log := util.Logger(cmd)
	client := util.Client(cmd)
	courseID := args[0]

	fmt.Println("modules:")
	for _, m := range client.ListModules(courseID) {
		log.Info().Msg(canvas.ModuleString(m))
	}
}
