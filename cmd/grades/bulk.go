package grades

import (
	"slices"
	"strconv"
	"time"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/modules"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"github.com/rs/zerolog/log"
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

var gradeStr string

func execGradesBulkCmd(cmd *cobra.Command, args []string) {
	log := util.Logger(cmd)
	client := util.Client(cmd)
	log.Info().Msg("Locating modules for bulk grading")

	grade, _ := cmd.Flags().GetFloat64("grade")
	gradeStr = strconv.FormatFloat(grade, 'f', 4, 64)

	// Make sure we are happy with the grade we are about to give
	if live, _ := cmd.Flags().GetBool("live"); live {
		if force, _ := cmd.Flags().GetBool("force"); !force {
			confirmGrade(cmd)
		}
	}

	// First select a module
	var module *canvasauto.Module
	if moduleID, _ := cmd.Flags().GetString("moduleID"); moduleID != "" {
		if m, err := client.GetModuleByID(util.GetCourseIdStr(cmd), moduleID); err != nil {
			log.Fatal().Err(err)
		} else if m == nil {
			log.Fatal().Msg("Module not found")
		} else {
			module = m
		}
	} else {
		module = util.MustFuzzyFindModule(cmd)
	}
	log.Info().
		Str("module", canvas.StrOrNil(module.Name)).
		Msg("Module selected")

	// Second select an assignment
	assignments := make([]*canvasauto.Assignment, 0)
	if allAssignments, _ := cmd.Flags().GetBool("allAssignments"); allAssignments {
		assignments = client.GetAssignmentsFromModule(util.GetCourseIdStr(cmd), module)
		if len(assignments) < 1 {
			log.Fatal().Str("id", canvas.StrOrNil(module.Id)).
				Str("name", canvas.StrOrNil(module.Name)).
				Msg("No assignments available for module")
		}
	} else {
		assignment := util.MustFuzzyFindAssignment(cmd, module)
		log.Info().
			Str("module", canvas.StrOrNil(module.Name)).
			Str("assignment", canvas.StrOrNil(assignment.Name)).
			Msg("Assignment selected")
		if assignment != nil {
			assignments = append(assignments, assignment)
		}
	}

	// Finally grade assignment(s)
	gradeAssignments(cmd, module, assignments)
}

func gradeAssignments(cmd *cobra.Command, module *canvasauto.Module, assignments []*canvasauto.Assignment) {
	for _, assignment := range assignments {
		gradeAssignment(cmd, module, assignment)
	}
}

func gradeAssignment(cmd *cobra.Command, module *canvasauto.Module, assignment *canvasauto.Assignment) {
	util.Logger(cmd).Info().
		Str("module", canvas.StrOrNil(module.Name)).
		Str("assignment", canvas.StrOrNil(assignment.Name)).
		Msg("Grading assignment submissions")
	submissions := getSubmissions(cmd, module, assignment)
	for _, s := range submissions {
		gradeSubmission(cmd, assignment, s)
	}
}

// For a given submission, check flags and record the provided grade
func gradeSubmission(cmd *cobra.Command, a *canvasauto.Assignment, s *canvasauto.Submission) {
	client := util.Client(cmd)

	user := client.GetUserById(util.GetCourseIdStr(cmd), *s.UserId)
	if user == nil {
		log.Debug().Int("userID", *s.UserId).Msg("Skipping unknown user")
		return
	}

	// Check if this is an interesting user
	emails, _ := cmd.Flags().GetStringArray("email")
	if len(emails) > 0 && !slices.Contains(emails, *user.Email) {
		log.Debug().Strs("emails", emails).Str("email", *user.Email).
			Msg("Skipping user not found in --email filter")
		return
	}

	// Unless forced to do so, don't overwrite work that is already submitted
	submitted, _ := cmd.Flags().GetBool("submitted")
	if !submitted && canvas.StrOrNil(s.WorkflowState) != "unsubmitted" {
		log.Debug().
			Str("student", canvas.StrOrNil(user.Email)).
			Msg("Skipping submitted assignment, set --submitted to grade")
		return
	}

	// Unless forced to do so, don't grade work that isn't late
	notLate, _ := cmd.Flags().GetBool("notLate")
	if !notLate && !time.Now().After(*a.DueAt) {
		log.Info().
			Str("student", canvas.StrOrNil(user.Email)).
			Time("dueAt", *a.DueAt).
			Msg("Skipping assignment that is not past due")
		return
	}

	log.Info().
		Str("state", canvas.StrOrNil(s.WorkflowState)).
		Str("curGrade", canvas.StrOrNil(s.Grade)).
		Str("due", canvas.StrOrNil(a.DueAt)).
		Str("student", canvas.StrOrNil(user.Email)).
		Msg("Found submission for bulk grading")

	if s.Grade != nil && cmpGradeStrings(*s.Grade, gradeStr) {
		log.Warn().
			Str("state", canvas.StrOrNil(s.WorkflowState)).
			Str("curGrade", canvas.StrOrNil(s.Grade)).
			Str("due", canvas.StrOrNil(a.DueAt)).
			Str("student", canvas.StrOrNil(user.Email)).
			Msg("Skipping item, has desired grade")
		return
	}

	if live, _ := cmd.Flags().GetBool("live"); !live {
		log.Debug().Msg("Not grading, set --live to grade")
		return
	}

	client.GradeSubmission(&canvas.UpdateSubmissionOpts{
		Score: gradeStr,
		ListSubmissionsOpts: &canvas.ListSubmissionsOpts{
			CourseID:   util.GetCourseIdStr(cmd),
			UserID:     canvas.StrOrNil(s.UserId),
			Assignment: a,
		},
	})
}

// Returns true if both strings converted to a float64
// evaluate to be identical
func cmpGradeStrings(s1 string, s2 string) bool {
	g1, _ := strconv.ParseFloat(s1, 64)
	g2, _ := strconv.ParseFloat(s2, 64)
	return g1 == g2
}

func confirmGrade(cmd *cobra.Command) {
	log.Warn().Str("grade", gradeStr).Msg("Applying this grade to all submissions!!! Hit ctrl+c in 10s to abort...")

	ticker := time.NewTicker(10 * time.Second)
	select {
	case <-cmd.Context().Done():
		log.Fatal().Msg("Aborted grading")
	case <-ticker.C:
		return
	}
}

func getSubmissions(cmd *cobra.Command, m *canvasauto.Module, a *canvasauto.Assignment) []*canvasauto.Submission {
	client := util.Client(cmd)
	submissions, err := client.ListAssignmentSubmissions(&canvas.ListSubmissionsOpts{
		CourseID:   util.GetCourseIdStr(cmd),
		Module:     m,
		Assignment: a,
	})
	if err != nil {
		log.Fatal().Err(err).
			Str("module", canvas.StrOrNil(m.Name)).
			Str("assignment", canvas.StrOrNil(a.Name)).
			Msg("Failed to list assignment submissions")
	}
	return submissions
}

func init() {
	gradesBulkCmd.Flags().StringP("moduleID", "m", "", "Specify a module ID")
	gradesBulkCmd.Flags().Bool("submitted", false, "CAUTION!! Will bulk mark submitted grades!")
	gradesBulkCmd.Flags().Bool("notLate", false, "CAUTION!! Will bulk mark grades that are not already late!")
	gradesBulkCmd.Flags().Bool("live", false, "CAUTION!! Will enable live grading!")
	gradesBulkCmd.Flags().Float64P("grade", "g", 0.0, "Grade to mark for assignment")
	gradesBulkCmd.Flags().BoolP("force", "f", false, "Don't warn first, just grade")
	gradesBulkCmd.Flags().BoolP("allAssignments", "a", false, "Grade all assignments in selected module")

	gradesBulkCmd.RegisterFlagCompletionFunc("moduleID", modules.ValidateModuleIdArg)
}
