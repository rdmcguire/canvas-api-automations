package modules

import (
	"strings"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"github.com/spf13/cobra"
)

func ValidateModuleIdArg(cmd *cobra.Command, args []string, toComplete string,
) ([]string, cobra.ShellCompDirective) {
	// First retrieve all modules
	client := util.Client(cmd)
	modules := client.ListModules(util.GetCourseIdStr(cmd))

	// Then filter and return
	validModules := make([]string, 0, len(modules))
	for _, m := range modules {
		if strings.HasPrefix(canvas.IntStrOrNil(m.Id), toComplete) {
			validModules = append(validModules, canvas.IntStrOrNil(m.Id))
		}
	}

	return validModules, cobra.ShellCompDirectiveNoFileComp
}
