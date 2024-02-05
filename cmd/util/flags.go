package util

import (
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func EnsureCourseId(cmd *cobra.Command, args []string) {
	var courseId int

	id, _ := cmd.Flags().GetInt("courseID")
	if id != 0 {
		return
	}

	if id, set := os.LookupEnv("CANVAS_COURSE_ID"); set {
		var err error
		if courseId, err = strconv.Atoi(id); err == nil && courseId != 0 {
			cmd.Flags().Set("courseID", id)
		}
	}

	if courseId == 0 {
		log.Fatal().
			Msg("CourseID is required, provide with --courseID or COURSE_ID env")
	}
}

func GetCourseIdStr(cmd *cobra.Command) string {
	id, _ := cmd.Flags().GetInt("courseID")
	return strconv.Itoa(id)
}

func GetCourseIdInt(cmd *cobra.Command) int {
	id, _ := cmd.Flags().GetInt("courseID")
	return id
}
