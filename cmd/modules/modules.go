package modules

import (
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"github.com/spf13/cobra"
)

var ModulesCmd = &cobra.Command{
	Use:              "modules",
	Aliases:          []string{"module", "m"},
	Short:            "Canvas modules",
	Long:             "Commands for interacting with modules in Canvas",
	PersistentPreRun: util.EnsureCourseId,
}

func init() {
	ModulesCmd.AddCommand(modulesShowCmd)
}
