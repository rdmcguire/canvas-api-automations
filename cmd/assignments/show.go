package assignments

import (
	"strconv"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"github.com/spf13/cobra"
)

var assignmentsShowCmd = &cobra.Command{
	Use:               "show",
	Aliases:           []string{"s", "ls"},
	Args:              cobra.NoArgs,
	ValidArgsFunction: cobra.NoFileCompletions,
	Short:             "Show assignments for a course",
	Run:               execAssignmentsShowCmd,
}

func execAssignmentsShowCmd(cmd *cobra.Command, args []string) {
	log := util.Logger(cmd)
	client := util.Client(cmd)

	courseID := util.GetCourseIdStr(cmd)
	moduleID, _ := cmd.Flags().GetInt("module")
	assignments := make([]*canvasauto.Assignment, 0)
	var err error

	if moduleID != 0 {
		log.Debug().Int("moduleID", moduleID).Msg("Listing assignments by module")
		module, err := client.GetModuleByID(courseID, strconv.Itoa(moduleID))
		if err != nil || module == nil {
			log.Fatal().Err(err).
				Int("courseID", util.GetCourseIdInt(cmd)).
				Int("moduleID", moduleID).Msg("Failed to find module")
		} else if module.Items == nil {
			log.Fatal().Msg("Module has no items")
		}
		for _, i := range *module.Items {
			if canvas.StrStrOrNil(i.Type) == "Assignment" {
				a, err := client.GetAssignmentById(&canvas.AssignmentOpts{
					ID:             canvas.IntStrOrNil(i.ContentId),
					ModuleItemOpts: &canvas.ModuleItemOpts{CourseID: courseID},
				})
				if err != nil {
					log.Error().Err(err).Msg("Failed retrieving assignment from module item")
					continue
				}
				assignments = append(assignments, a)
			}
		}
	} else {
		assignments, err = client.ListAssignments(courseID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to list assignments")
		}
	}

	for _, a := range assignments {
		log.Info().Msg(canvas.AssignmentString(a))
	}
}

func init() {
	assignmentsShowCmd.Flags().Int("module", 0, "Specify module by ID")
}
