package canvas

import (
	"encoding/json"
	"strings"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"github.com/rs/zerolog/log"
)

func (c *Client) GetUserByEmail(courseID string, email string) *canvasauto.User {
	for _, user := range c.ListUsersInCourse(courseID, "") {
		if strings.ToLower(StrOrNil(user.Email)) == strings.ToLower(email) {
			return user
		}
	}
	return nil
}

func (c *Client) ListUsersInCourse(courseID string, search string) []*canvasauto.User {
	users := make([]*canvasauto.User, 0)
	role := "StudentEnrollment"
	page := 1
	for {
		pageUsers := make([]*canvasauto.User, 0)
		params := &canvasauto.ListUsersInCourseUsersParams{
			Page:           &page,
			PerPage:        &DefPerPage,
			EnrollmentRole: &role,
		}
		if search != "" {
			params.SearchTerm = &search
		}

		r, err := c.api.ListUsersInCourseUsers(c.ctx, courseID, params)
		if err != nil || r.Body == nil {
			log.Error().Err(err)
		}

		json.NewDecoder(r.Body).Decode(&pageUsers)
		users = append(users, pageUsers...)
		if isLastPage(r) {
			break
		}
		page += 1
	}
	return users
}
