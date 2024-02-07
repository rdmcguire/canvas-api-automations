package students

import (
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var studentsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"show", "ls", "l", "s"},
	Short:   "List students in the course",
	Run:     execStudentsListCmd,
}

func execStudentsListCmd(cmd *cobra.Command, args []string) {
	client := util.Client(cmd)
	courseID := util.GetCourseIdStr(cmd)

	util.Logger(cmd).Debug().Str("courseID", courseID).Msg("Listing students")
	students := client.ListStudentsInCourse(courseID)
	for _, student := range students {
		log.Info().
			Str("Name", canvas.StrOrNil(student.Name)).
			Str("Email", canvas.StrOrNil(student.Email)).
			Int64("ID", student.Id).
			Msg("")
	}
}
