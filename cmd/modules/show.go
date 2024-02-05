package modules

import (
	"fmt"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"github.com/spf13/cobra"
)

var modulesShowCmd = &cobra.Command{
	Use:               "show",
	Aliases:           []string{"ls", "s"},
	Args:              cobra.NoArgs,
	ValidArgsFunction: cobra.NoFileCompletions,
	Run:               execModulesShowCmd,
	Short:             "Show modules",
}

func execModulesShowCmd(cmd *cobra.Command, args []string) {
	log := util.Logger(cmd)
	client := util.Client(cmd)

	fmt.Println("modules:")
	for _, m := range client.ListModules(util.GetCourseIdStr(cmd)) {
		log.Info().Msg(canvas.ModuleString(m))
	}
}
