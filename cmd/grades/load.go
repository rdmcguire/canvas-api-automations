package grades

import (
	"fmt"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/netacad"
	"github.com/spf13/cobra"
)

var gradesLoadCmd = &cobra.Command{
	Use:     "load",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"import"},
	Short:   "Load and dump grades",
	Long:    `Mosly useful for debugging export from Netacad`,
	Run:     execGradesLoadCmd,
}

func execGradesLoadCmd(cmd *cobra.Command, args []string) {
	log := util.Logger(cmd)

	grades, err := netacad.LoadGradesFromFile(args[0])
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load grades export")
	}

	for student, grades := range *grades {
		if grades != nil {
			fmt.Println(student)
			for item, grade := range *grades {
				if grade != nil {
					fmt.Printf("\tItem %s, %+v\n", item, grade)
				}
			}
		}
	}
}

func init() {
}
