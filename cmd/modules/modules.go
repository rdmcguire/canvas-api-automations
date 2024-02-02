package modules

import (
	"github.com/spf13/cobra"
)

var ModulesCmd = &cobra.Command{
	Use:     "modules",
	Aliases: []string{"module", "m"},
	Short:   "Canvas modules",
	Long:    "Commands for interacting with modules in Canvas",
}

func init() {
	ModulesCmd.AddCommand(modulesShowCmd)
}
