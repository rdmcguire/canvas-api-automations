package grades

import (
	"fmt"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

var gradesBulkCmd = &cobra.Command{
	Use:     "bulk",
	Aliases: []string{"mark", "set"},
	Short:   "Marks grades, typically used to bulk mark late work to zero",
	Long: `This command will dynamically allow you to select the module, and
the assignment within the module, and confirm that you will be marking all
unsubmitted work with a grade of 0, or a grade provided by the
--grade / -g flag`,
	Run: execGradesBulkCmd,
}

func execGradesBulkCmd(cmd *cobra.Command, args []string) {
	log := util.Logger(cmd)

	log.Info().Msg("Locating modules for bulk grading")

	module := FuzzyFindModule(cmd)

	log.Info().Str("module", canvas.StrOrNil(module.Name)).
		Msg("Module selected, now selecting an assignment")

	assignment := FuzzyFindAssignment(cmd, module)
	log.Info().Str("module", canvas.StrOrNil(module.Name)).
		Str("assignment", canvas.StrOrNil(assignment.Name)).
		Msg("Assignment selected, locating submissions")
}

func FuzzyFindAssignment(cmd *cobra.Command, module *canvasauto.Module) *canvasauto.Assignment {
	log := util.Logger(cmd)
	client := util.Client(cmd)

	assignments := client.GetAssignmentsFromModule(util.GetCourseIdStr(cmd), module)
	if len(assignments) < 1 {
		log.Fatal().Str("id", canvas.StrOrNil(module.Id)).
			Str("name", canvas.StrOrNil(module.Name)).
			Msg("No assignments available for module")
	}

	idx, err := fuzzyfinder.Find(assignments, func(i int) string {
		return fmt.Sprintf("%s - %s - %s",
			canvas.StrOrNil(module.Name),
			canvas.StrOrNil(assignments[i].Id),
			canvas.StrOrNil(assignments[i].Name))
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed selecting assignment from module")
	}

	return assignments[idx]
}

func FuzzyFindModule(cmd *cobra.Command) *canvasauto.Module {
	log := util.Logger(cmd)
	client := util.Client(cmd)

	modules := client.ListModules(util.GetCourseIdStr(cmd))
	idx, err := fuzzyfinder.Find(modules, func(i int) string {
		return fmt.Sprintf("%s - %s",
			canvas.StrOrNil(modules[i].Id),
			canvas.StrOrNil(modules[i].Name))
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to select a module")
	}
	return modules[idx]
}

func init() {
	gradesBulkCmd.Flags().Bool("submitted", false, "CAUTION!! Will bulk mark submitted grades!")
	gradesBulkCmd.Flags().Float64P("grade", "g", 0.0, "Grade to mark for assignment")
}
