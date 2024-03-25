package grades

import (
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/netacad"
	"github.com/spf13/cobra"
)

var gradesFindCmd = &cobra.Command{
	Use:     "find <grade_export.csv>",
	Aliases: []string{"match", "matches", "f"},
	Args:    cobra.ExactArgs(1),
	Run:     execGradesFindCmd,
	Short:   "Loads grade export and finds matching course assignments",
	Long: `This command will load a grade export from Netacad,
then locate matching courses in the Canvas API for the current course.
This is intended to be for debugging the assignment matching logic.`,
}

func execGradesFindCmd(cmd *cobra.Command, args []string) {
	log := util.Logger(cmd)

	assignments, err := netacad.Assignments(&netacad.LoadGradesFromFileOpts{
		File: args[0],
	})

	if err != nil {
		log.Fatal().Err(err).
			Msg("Failed finding assignments in export, do not push!")
	}

	log.Info().Int("assignments", len(assignments)).
		Msg("Netacad Grade Export Loaded")

	for _, a := range assignments {
		log.Debug().Str("item", a).Msg("Locating assignment")
		assignment, module := getAssignmentFromGrade(cmd, a)

		if assignment != nil {
			log.Info().Str("item", a).
				Str("module", canvas.StrOrNil(module.Name)).
				Str("canvasItem", canvas.StrOrNil(assignment.Name)).
				Msg("Found canvas assignment")
		}
	}
}
