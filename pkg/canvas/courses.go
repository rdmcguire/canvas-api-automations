package canvas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tomnomnom/linkheader"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
)

var DefPerPage = 20

type CourseResponse struct {
	Course *canvasauto.Course `json:"course"`
}

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

func (c *Client) ListStudentsInCourse(id string) []*canvasauto.User {
	users := make([]*canvasauto.User, 0)
	role := "StudentEnrollment"
	page := 1
	for {
		pageUsers := make([]*canvasauto.User, 0)
		r, _ := c.api.ListUsersInCourseUsers(c.ctx, id,
			&canvasauto.ListUsersInCourseUsersParams{
				Page:           &page,
				PerPage:        &DefPerPage,
				EnrollmentRole: &role,
			},
		)
		json.NewDecoder(r.Body).Decode(&pageUsers)
		users = append(users, pageUsers...)
		if isLastPage(r) {
			fmt.Println("Last Page")
			break
		}
		page += 1
	}
	return users
}

func isLastPage(r *http.Response) bool {
	links := linkheader.Parse(r.Header.Get("link"))
	for _, link := range links {
		if link.Rel == "next" {
			return false
		}
	}
	return true
}