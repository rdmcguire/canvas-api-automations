package grades

import (
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"github.com/spf13/cobra"
)

var gradesPushCmd = &cobra.Command{
	Use:     "push <grade_export.csv>",
	Aliases: []string{"sync", "update"},
	Args:    cobra.ExactArgs(1),
	Run:     execGradesPushCmd,
	Short:   "Pushes exported grades into Canvas",
	Long: `This command will attempt to set any grades
in Canvas that are present in the export. If you want to
replace any existing grades, add --overwrite. For safety,
you must add --live to commit the grades, else it will simply
display proposed changes to the gradebook.`,
}

func execGradesPushCmd(cmd *cobra.Command, args []string) {
	log := util.Logger(cmd)

	gradebook := mustLoadGrades(cmd, args[0])

	for student, grades := range *gradebook {
		if grades.Count() == 0 {
			continue
		}

		log.Info().Str("email", student.Email).
			Int("gradesLoaded", grades.Count()).
			Str("name", student.First+" "+student.Last).
			Msg("Grading Student")
		log.Debug().Str("email", student.Email).
			Any("gradesLoaded", *grades).
			Msg("Student Grades Loaded")

		gradeStudent(cmd, &student, grades)
	}
}

func init() {
	gradesPushCmd.Flags().Bool("overwrite", false,
		"Set to overwrite existing grades, otherwise skips them")
	gradesPushCmd.Flags().Bool("live", false,
		"Set to actually push grades into Canvas, otherwise reports changes")
}
