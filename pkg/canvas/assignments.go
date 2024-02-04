package canvas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"k8s.io/utils/ptr"
)

type AssignmentOpts struct {
	ID         string
	Assignment *canvasauto.Assignment
	*ModuleItemOpts
}

func (c *Client) UpdateAssignment(opts *AssignmentOpts) (*canvasauto.Assignment, error) {
	r, err := c.api.EditAssignment(c.ctx, opts.CourseID, opts.ID,
		canvasauto.EditAssignmentJSONRequestBody{
			AssignmentDescription: opts.Assignment.Description,
		})
	if err != nil {
		return nil, err
	}

	return decodeAssignmentResponse(r)
}

func (c *Client) GetAssignmentById(opts *AssignmentOpts) (*canvasauto.Assignment, error) {
	r, err := c.api.GetSingleAssignment(c.ctx,
		opts.CourseID,
		opts.ID,
		&canvasauto.GetSingleAssignmentParams{})
	if err != nil {
		return nil, err
	}

	return decodeAssignmentResponse(r)
}

func (c *Client) ListAssignments(courseID string, modules ...int) ([]*canvasauto.Assignment, error) {
	assignments := make([]*canvasauto.Assignment, 0)
	opts := &canvasauto.ListAssignmentsParams{Page: ptr.To("1")}
	page := 1
	for {
		pageAssignments := make([]*canvasauto.Assignment, 0)
		r, err := c.api.ListAssignments(c.ctx, courseID, opts)
		if err != nil {
			return nil, err
		}

		json.NewDecoder(r.Body).Decode(&pageAssignments)
		assignments = append(assignments, pageAssignments...)

		if isLastPage(r) {
			break
		}

		page += 1
		opts.Page = ptr.To(strconv.Itoa(page))
	}
	return assignments, nil
}

func decodeAssignmentResponse(r *http.Response) (*canvasauto.Assignment, error) {
	assignment := &canvasauto.Assignment{}
	decoder := json.NewDecoder(r.Body)
	return assignment, decoder.Decode(assignment)
}

func AssignmentString(assignment *canvasauto.Assignment) string {
	return fmt.Sprintf("Assignment %s [ID:%s] [Published:%s]",
		StrStrOrNil(assignment.Name),
		IntStrOrNil(assignment.Id),
		BoolStrOrNil(assignment.Published),
	)
}
