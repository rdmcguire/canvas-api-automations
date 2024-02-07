package grades

import (
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/netacad"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func gradeStudent(cmd *cobra.Command, student *netacad.Student, grades *netacad.Grades) {
	for item, _ := range *grades {
		log.Debug().Str("student", student.Email).
			Str("item", item).
			Msg("Locating item for student grade")
		getAssignmentFromGrade(cmd, item)
	}
}

func getAssignmentFromGrade(cmd *cobra.Command, item string) *canvasauto.Assignment {
	var assignment *canvasauto.Assignment
	client := util.Client(cmd)

	getOpts := &canvas.ModuleItemOpts{
		CourseID:    util.GetCourseIdStr(cmd),
		Name:        item,
		Insensitive: true,
	}

	// Try with a full match (probably a waste of time)
	if found := client.FindItem(getOpts); found != nil {
		getAssignmentFromItem(cmd, found, getOpts)
	}

	return assignment
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
