package canvas

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
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
