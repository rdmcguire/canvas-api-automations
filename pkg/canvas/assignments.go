package canvas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"golang.org/x/exp/slices"
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

func (c *Client) ListAssignmentsByModule(courseID string, moduleIDs ...int,
) map[*canvasauto.Module][]*canvasauto.Assignment {
	moduleAssignments := make(map[*canvasauto.Module][]*canvasauto.Assignment, 0)

	// First get modules
	modules := c.ListModules(courseID)
	for _, module := range modules {
		if len(moduleIDs) > 0 && !slices.Contains(moduleIDs, *module.Id) {
			continue
		}
		moduleAssignments[module] = c.GetAssignmentsFromModule(courseID, module)
	}
	return moduleAssignments
}

func (c *Client) GetAssignmentsFromModule(courseID string, module *canvasauto.Module) []*canvasauto.Assignment {
	assignments := make([]*canvasauto.Assignment, 0)
	if module.Items == nil {
		return assignments
	}

	for _, i := range *module.Items {
		if StrOrNil(i.Type) == "Assignment" {
			assignment, err := c.GetAssignmentById(&AssignmentOpts{
				ModuleItemOpts: &ModuleItemOpts{
					CourseID: courseID,
				},
				ID: StrOrNil(i.ContentId),
			})
			if err == nil && assignment != nil {
				assignments = append(assignments, assignment)
			}
		}
	}
	return assignments
}

func (c *Client) ListAssignments(courseID string) ([]*canvasauto.Assignment, error) {
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
	return fmt.Sprintf("Assignment %s [ID:%s] [Published:%s] [Due:%s]",
		StrOrNil(assignment.Name),
		StrOrNil(assignment.Id),
		StrOrNil(assignment.Published),
		StrOrNil(assignment.DueAt),
	)
}
