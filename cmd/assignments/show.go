package assignments

import (
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/courses"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"github.com/spf13/cobra"
)

var assignmentsShowCmd = &cobra.Command{
	Use:               "show (courseID)",
	Aliases:           []string{"s", "ls"},
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: courses.ValidateCourseIdArg,
	Short:             "Show assignments for a course",
	Run:               execAssignmentsShowCmd,
}

func execAssignmentsShowCmd(cmd *cobra.Command, args []string) {
	log := util.Logger(cmd)
	client := util.Client(cmd)

	assignments, err := client.ListAssignments(args[0])
	if err != nil {
		log.Error().Err(err).Msg("Failed to list assignments")
	}

	for _, a := range assignments {
		log.Info().Msg(canvas.AssignmentString(a))
	}
}
