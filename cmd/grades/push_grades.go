package grades

import (
	"fmt"
	"regexp"
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
	externalToolRegexp = regexp.MustCompile(`External tool: (.*)`)
)

func gradeStudent(cmd *cobra.Command, student *netacad.Student, grades *netacad.Grades) {
	getAssignmentsFromGrades(cmd, student, grades) // First find matching assignments
	updateGrades(cmd, student, grades)             // Then update grades
}

func updateGrades(cmd *cobra.Command, student *netacad.Student, grades *netacad.Grades) {
	for _, grade := range *grades {
		if grade.Assignment != nil {
			updateGrade(cmd, student, grade)
		}
	}
}

func updateGrade(cmd *cobra.Command, student *netacad.Student, grade *netacad.Grade) {
	live, _ := cmd.Flags().GetBool("live")
	if !live {
		util.Logger(cmd).Info().
			Str("student", student.Email).
			Str("module", canvas.StrOrNil(grade.Module.Name)).
			Str("assignment", canvas.StrOrNil(grade.Assignment.Name)).
			Float64("grade", grade.Percentage).
			Msg("DRY RUN grade update")
	}
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

Found:
	if assignment != nil {
		log.Debug().
			Str("assignment", canvas.StrOrNil(assignment.Name)).
			Str("item", getOpts.Name).
			Str("Module", canvas.StrOrNil(getOpts.Module.Name)).
			Msg("Match found!")
	} else {
		log.Warn().Str("item", getOpts.Name).Msg("Failed to locate grade item match")
	}

	return assignment, getOpts.Module
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
	log.Debug().Str("name", newOpts.Name).Msg("Attempting exam -> quiz withour leading zeroes")
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
