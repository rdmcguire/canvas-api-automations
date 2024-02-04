package modules

import (
	"fmt"
	"strconv"

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
	courseID, _ := cmd.Flags().GetInt("courseID")

	fmt.Println("modules:")
	for _, m := range client.ListModules(strconv.Itoa(courseID)) {
		log.Info().Msg(canvas.ModuleString(m))
	}
}
