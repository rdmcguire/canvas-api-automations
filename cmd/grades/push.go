package grades

import (
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/netacad"
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
	emails, err := cmd.Flags().GetStringArray("email")
	if err != nil {
		log.Error().Err(err).Send()
	}

	if len(emails) > 0 {
		log.Info().Strs("emails", emails).
			Msg("Student email filter loaded")
	}

	// Grade students concurrently
	gradeWg := &sync.WaitGroup{}
	for student, grades := range *gradebook {
		if grades.Count() == 0 {
			log.Info().Str("email", student.Email).
				Msg("Student has nothing to grade")
			continue
		} else if !studentInFilter(emails, &student) {
			log.Debug().Str("email", student.Email).
				Msg("Skipping student by filter")
			continue
		}

		gradeWg.Add(1)
		go func(g *netacad.Grades, s netacad.Student) {
			defer gradeWg.Done()
			grade(cmd, g, s)
		}(grades, student)
	}

	log.Debug().Msg("All grading goroutines launched, waiting for completion")
	gradeWg.Wait()
}

func grade(cmd *cobra.Command, grades *netacad.Grades, student netacad.Student) {
	log.Info().Str("email", student.Email).
		Int("gradesLoaded", grades.Count()).
		Str("name", student.First+" "+student.Last).
		Msg("Launched Grading Goroutine")
	log.Debug().Str("email", student.Email).
		Any("gradesLoaded", *grades).
		Msg("Student Grades Loaded")
	gradeStudent(cmd, &student, grades)
}

func studentInFilter(filter []string, student *netacad.Student) bool {
	if len(filter) > 0 {
		for _, email := range filter {
			if strings.EqualFold(email, student.Email) {
				return true
			}
		}
	} else if len(filter) == 0 {
		return true
	}

	return false
}

func init() {
	gradesPushCmd.Flags().Bool("overwrite", false,
		"Set to overwrite existing grades, otherwise skips them")
	gradesPushCmd.Flags().Bool("live", false,
		"Set to actually push grades into Canvas, otherwise reports changes")
}
