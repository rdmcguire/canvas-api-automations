package grades

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/netacad"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Patterns for grade items
var (
	chapterExamRegexp  = regexp.MustCompile(`.*Chapter (\d+) Exam`)
	chapterLabRegexp   = regexp.MustCompile(`.*Lab (\d+)`)
	externalToolRegexp = regexp.MustCompile(`External tool: (.*)`)
	finalRegexp        = regexp.MustCompile(`(Final.*Exam)`)
)

func gradeStudent(cmd *cobra.Command, student *netacad.Student, grades *netacad.Grades) {
	getAssignmentsFromGrades(cmd, student, grades) // First find matching assignments
	getSubmissionsFromGrades(cmd, student, grades)
	updateGrades(cmd, student, grades) // Then update grades
}

func updateGrades(cmd *cobra.Command, student *netacad.Student, grades *netacad.Grades) {
	for _, grade := range *grades {
		if grade.Assignment != nil {
			updateGrade(cmd, student, grade)
		}
	}
}

func updateGrade(cmd *cobra.Command, student *netacad.Student, grade *netacad.Grade) {
	log := util.Logger(cmd)

	live, _ := cmd.Flags().GetBool("live")
	if !live {
		log.Info().
			Str("student", student.Email).
			Str("module", canvas.StrOrNil(grade.Module.Name)).
			Str("assignment", canvas.StrOrNil(grade.Assignment.Name)).
			Int("submissions", len(grade.Submissions)).
			Float64("percentage", grade.Percentage).
			Float64("points", grade.Grade).
			Str("pointsPossible", canvas.StrOrNil(grade.Assignment.PointsPossible)).
			Msg("DRY RUN grade update")
	} else {
		log.Warn().
			Str("student", student.Email).
			Str("assignment", canvas.StrOrNil(grade.Assignment.Name)).
			Str("score", fmt.Sprintf("%.3f%%", grade.Percentage)).
			Msg("Performing live grading")
		util.Client(cmd).GradeSubmission(&canvas.UpdateSubmissionOpts{
			Score: fmt.Sprintf("%.3f%%", grade.Percentage),
			ListSubmissionsOpts: &canvas.ListSubmissionsOpts{
				CourseID:   util.GetCourseIdStr(cmd),
				UserID:     strconv.Itoa(int(grade.User.Id)),
				Assignment: grade.Assignment,
			},
		})
	}
}

func getSubmissionsFromGrades(cmd *cobra.Command, student *netacad.Student, grades *netacad.Grades) {
	for _, grade := range *grades {
		if grade.Assignment == nil {
			continue
		}

		getSubmissionsForGrade(cmd, student, grade)
	}
}

func getSubmissionsForGrade(cmd *cobra.Command, student *netacad.Student, grade *netacad.Grade) {
	client := util.Client(cmd)
	log := util.Logger(cmd)
	overwrite, _ := cmd.Flags().GetBool("overwrite")

	user := client.GetUserByEmail(util.GetCourseIdStr(cmd), student.Email)
	if user == nil {
		log.Error().Str("student", student.Email).Msg("Failed to locate user")
	}

	grade.User = user

	submissions, err := client.ListAssignmentSubmissions(&canvas.ListSubmissionsOpts{
		CourseID:   util.GetCourseIdStr(cmd),
		Assignment: grade.Assignment,
		UserID:     strconv.Itoa(int(user.Id)), // filter specific user ID
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed listing submissions for user")
		return
	}

	var submission *canvasauto.Submission
	if len(submissions) > 0 {
		submission = submissions[len(submissions)-1]
	}

	score := ScaleGradeToAssignment(grade.Percentage, *grade.Assignment.PointsPossible)
	log.Info().Msg(fmt.Sprintf("%s scored %.2f/%.2f (originally %.2f[%.2f%%]) on %s",
		student.Email,
		score, *grade.Assignment.PointsPossible,
		grade.Grade, grade.Percentage,
		*grade.Assignment.Name))

	if canvas.StrOrNil(submission.WorkflowState) != "unsubmitted" {
		if overwrite {
			log.Warn().
				Str("student", student.Email).
				Str("assignment", *grade.Assignment.Name).
				Msg("Grade already submitted but overwrite enabled, grading forced!")
		} else {
			log.Info().
				Str("student", student.Email).
				Str("assignment", *grade.Assignment.Name).
				Msg("Grade already submitted, skipping...")
			grade.Assignment = nil
		}
	}
}

func ScaleGradeToAssignment(gradePercent float64, assignmentPointsPossible float32) float32 {
	return (float32(gradePercent) / 100) * assignmentPointsPossible
}

func getAssignmentsFromGrades(cmd *cobra.Command, student *netacad.Student, grades *netacad.Grades) {
	for item, grade := range *grades {
		log.Debug().Str("student", student.Email).
			Str("item", item).
			Msg("Locating item for student grade")
		grade.Assignment, grade.Module = getAssignmentFromGrade(cmd, item)
	}
}

func getAssignmentFromGrade(cmd *cobra.Command, item string,
) (*canvasauto.Assignment, *canvasauto.Module) {
	var assignment *canvasauto.Assignment

	// Skip total columns
	if strings.HasSuffix(strings.ToLower(item), "total") {
		log.Debug().Str("item", item).Msg("Skipping total item")
		return nil, nil
	}

	getOpts := &canvas.ModuleItemOpts{
		CourseID:    util.GetCourseIdStr(cmd),
		Name:        item,
		Insensitive: true,
		Fuzzy:       false,
	}

	if assignmentCache == nil {
		assignmentCache = util.NewAssignmentCache()
	}

	// First check if we have already seen this and failed to find it
	if assignmentCache.LostCause(item) {
		log.Debug().Str("item", item).Msg("Skipping item, already known to be a lost cause")
		return nil, nil
	}

	// Then try to retrieve from cache
	if assignment, getOpts.Module = assignmentCache.Get(item); assignment != nil {
		log.Debug().
			Str("name", *assignment.Name).
			Msg("Found assignment in assignment cache")
		return assignment, getOpts.Module
	}

	// Try with a full match (probably a waste of time)
	if assignment = tryFullNameMatch(cmd, getOpts); assignment != nil {
		goto Found
	}

	// Try without "External tool: " prefix if set
	if assignment = tryExternalToolMatch(cmd, getOpts); assignment != nil {
		goto Found
	}

	// Try extracting from "Chapter XX Exam" -> "Quiz XX"
	if assignment = tryExamToQuiz(cmd, getOpts); assignment != nil {
		goto Found
	}

	// Try extracting from "Lab 0X" -> "Lab X"
	if assignment = tryLabZeroPadded(cmd, getOpts); assignment != nil {
		goto Found
	}

	// Try looking for a midterm assignment
	if assignment = tryMidtermExam(cmd, getOpts); assignment != nil {
		goto Found
	}

	// Try looking for a final assignment
	if assignment = tryFinalExam(cmd, getOpts); assignment != nil {
		goto Found
	}

Found:
	// Add found item to cache and return it
	if assignment != nil {
		assignmentCache.Set(item, assignment, getOpts.Module)
		log.Debug().
			Str("assignment", canvas.StrOrNil(assignment.Name)).
			Str("item", getOpts.Name).
			Str("Module", canvas.StrOrNil(getOpts.Module.Name)).
			Msg("Match found!")
	} else {
		log.Warn().Str("item", getOpts.Name).Msg("Failed to locate grade item match")
		assignmentCache.IsLostCause(item) // Don't bother trying a second time
	}

	return assignment, getOpts.Module
}

// Fuzzy finds either Final Exam or Final Comprehensive Exam
func tryFinalExam(cmd *cobra.Command, opts *canvas.ModuleItemOpts) *canvasauto.Assignment {
	matches := finalRegexp.FindStringSubmatch(opts.Name)
	if len(matches) != 2 {
		return nil
	}

	newOpts := *opts
	newOpts.Name = matches[1]

	// This is necessary as the loaded courses in canvas doesn't
	// necessarily contain the correct chapters and, frankly it doesn't matter.
	// There is only one final.
	newOpts.Contains = true

	log.Debug().Str("name", newOpts.Name).Msg("Attempting final exam match")
	if found := util.Client(cmd).FindItem(&newOpts); found != nil {
		opts.Module = getModuleFromItem(cmd, found)
		return getAssignmentFromItem(cmd, found, &newOpts)
	}
	return nil
}

// Attempts to locate an assignment where the name is Midterm Exam
// but the exported item ends with (Modules x-x)
func tryMidtermExam(cmd *cobra.Command, opts *canvas.ModuleItemOpts) *canvasauto.Assignment {
	if !strings.Contains(opts.Name, "Midterm Exam") {
		return nil
	}
	newOpts := *opts

	newOpts.Name = "Midterm Exam"
	log.Debug().Str("name", newOpts.Name).Msg("Checking for Midterm Exam")
	if found := util.Client(cmd).FindItem(&newOpts); found != nil {
		opts.Module = getModuleFromItem(cmd, found)
		return getAssignmentFromItem(cmd, found, &newOpts)
	}

	return nil
}

// Attempts to return an item with the original grade item
// Also, what the FUCK Netacad. Everything is obnoxiously inconsistent.
func tryFullNameMatch(cmd *cobra.Command, opts *canvas.ModuleItemOpts) *canvasauto.Assignment {
	log.Debug().Str("name", opts.Name).Msg("Attempting full match")
	if found := util.Client(cmd).FindItem(opts); found != nil {
		opts.Module = getModuleFromItem(cmd, found)
		return getAssignmentFromItem(cmd, found, opts)
	}
	return nil
}

// Attempts to locate an assignment where the name is Lab X but
// the graded item is zero-padded. Also strips the annoying
// External Tool: prefix
func tryLabZeroPadded(cmd *cobra.Command, opts *canvas.ModuleItemOpts) *canvasauto.Assignment {
	matches := chapterLabRegexp.FindStringSubmatch(opts.Name)
	if len(matches) < 2 {
		return nil
	}

	newOpts := *opts

	// Try stripping leading zeroes
	newOpts.Name = fmt.Sprintf("Lab %s", strings.TrimPrefix(matches[1], "0"))
	log.Debug().Str("name", newOpts.Name).Msg("Attempting lab without leading zeroes")
	if found := util.Client(cmd).FindItem(&newOpts); found != nil {
		opts.Module = getModuleFromItem(cmd, found)
		return getAssignmentFromItem(cmd, found, &newOpts)
	}

	return nil
}

// Strips "External tool: " from the name
func tryExternalToolMatch(cmd *cobra.Command, opts *canvas.ModuleItemOpts) *canvasauto.Assignment {
	matches := externalToolRegexp.FindStringSubmatch(opts.Name)
	if len(matches) != 2 {
		return nil
	}

	newOpts := *opts
	newOpts.Name = matches[1]
	log.Debug().Str("name", newOpts.Name).Msg("Attempting strip external tool")
	if found := util.Client(cmd).FindItem(&newOpts); found != nil {
		opts.Module = getModuleFromItem(cmd, found)
		return getAssignmentFromItem(cmd, found, &newOpts)
	}
	return nil
}

// Converts Chapter XX Exam -> Quiz XX
// Also try without leading zeroes
func tryExamToQuiz(cmd *cobra.Command, opts *canvas.ModuleItemOpts) *canvasauto.Assignment {
	matches := chapterExamRegexp.FindStringSubmatch(opts.Name)
	if len(matches) < 2 {
		return nil
	}

	// Try with original number
	newOpts := *opts
	newOpts.Name = fmt.Sprintf("Quiz %s", matches[1])
	log.Debug().Str("name", newOpts.Name).Msg("Attempting exam -> quiz")
	if found := util.Client(cmd).FindItem(&newOpts); found != nil {
		opts.Module = getModuleFromItem(cmd, found)
		return getAssignmentFromItem(cmd, found, &newOpts)
	}

	// Try stripping leading zeroes
	newOpts.Name = fmt.Sprintf("Quiz %s", strings.TrimPrefix(matches[1], "0"))
	log.Debug().Str("name", newOpts.Name).Msg("Attempting exam -> quiz without leading zeroes")
	if found := util.Client(cmd).FindItem(&newOpts); found != nil {
		opts.Module = getModuleFromItem(cmd, found)
		return getAssignmentFromItem(cmd, found, &newOpts)
	}

	return nil
}

func getModuleFromItem(cmd *cobra.Command, item *canvasauto.ModuleItem) *canvasauto.Module {
	var module *canvasauto.Module
	module, _ = util.Client(cmd).
		GetModuleByID(util.GetCourseIdStr(cmd), canvas.StrOrNil(item.ModuleId))
	return module
}

func getAssignmentFromItem(cmd *cobra.Command, item *canvasauto.ModuleItem, opts *canvas.ModuleItemOpts) *canvasauto.Assignment {
	client := util.Client(cmd)
	var assignment *canvasauto.Assignment

	assignment, err := client.GetAssignmentById(&canvas.AssignmentOpts{
		ID: canvas.StrOrNil(item.ContentId),
		ModuleItemOpts: &canvas.ModuleItemOpts{
			CourseID: util.GetCourseIdStr(cmd),
		},
	})
	if err != nil {
		util.Logger(cmd).Error().Err(err)
		return nil
	}

	return assignment
}
