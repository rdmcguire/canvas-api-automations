package canvas

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
)

var DefPerPage = 20

func (c *Client) ListCourses() []*canvasauto.Course {
	courses := make([]*canvasauto.Course, 0)
	page := 1
	for {
		pageCourses := make([]*canvasauto.Course, 0)
		pageStr := strconv.Itoa(page)
		r, _ := c.api.ListYourCourses(c.ctx,
			&canvasauto.ListYourCoursesParams{
				PerPage: &DefPerPage,
				Page:    &pageStr,
			},
		)
		json.NewDecoder(r.Body).Decode(&pageCourses)
		courses = append(courses, pageCourses...)
		if isLastPage(r) {
			break
		}
		page += 1
	}
	return courses
}

func CourseString(course *canvasauto.Course) string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("ID:%d %s [%s]",
		*course.Id,
		*course.Name,
		*course.WorkflowState,
	))
	if course.StartAt != nil {
		str.WriteString(fmt.Sprintf(" @%s", *course.StartAt))
	}
	return str.String()
}
