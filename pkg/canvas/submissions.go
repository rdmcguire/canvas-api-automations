package canvas

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"k8s.io/utils/ptr"
)

const gradeMsg = "graded by automation <https://github.com/rdmcguire/canvas-api-automations>"

type UpdateSubmissionOpts struct {
	// Score can be in multiple formats
	// Reference https://canvas.instructure.com/doc/api/submissions.html#method.submissions_api.update
	Score string
	*ListSubmissionsOpts
}

type ListSubmissionsOpts struct {
	CourseID   string
	UserID     string
	Module     *canvasauto.Module
	Assignment *canvasauto.Assignment
}

func (c *Client) GradeSubmission(opts *UpdateSubmissionOpts) error {
	update := &canvasauto.GradeOrCommentOnSubmissionSectionsJSONBody{
		SubmissionPostedGrade: ptr.To(opts.Score),
		CommentTextComment:    ptr.To(gradeMsg),
	}
	updateBody, _ := json.Marshal(*update)

	r, err := c.api.GradeOrCommentOnSubmissionCoursesWithBody(c.ctx,
		opts.CourseID,
		StrOrNil(opts.Assignment.Id),
		opts.UserID,
		"application/json",
		bytes.NewReader(updateBody),
	)
	if r != nil && r.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Received non-200 response for submission update (%d)",
			r.StatusCode))
	}

	return err
}

func (c *Client) ListAssignmentSubmissions(opts *ListSubmissionsOpts) ([]*canvasauto.Submission, error) {
	var err error
	aID := strconv.Itoa(*opts.Assignment.Id)
	submissions := make([]*canvasauto.Submission, 0)

	page := 1
	listOpts := &canvasauto.ListAssignmentSubmissionsCoursesParams{
		Include: ptr.To([]string{"assignment", "user", "submission_history"}),
		Page:    &page,
	}
	for {
		var r *http.Response
		r, err = c.api.ListAssignmentSubmissionsCourses(c.ctx, opts.CourseID, aID, listOpts)
		if err != nil {
			break
		}

		pageSubmissions := make([]*canvasauto.Submission, 0)
		json.NewDecoder(r.Body).Decode(&pageSubmissions)

		if opts.UserID != "" {
			for _, submission := range pageSubmissions {
				if StrOrNil(submission.UserId) == opts.UserID {
					submissions = append(submissions, submission)
				}
			}
		} else {
			submissions = append(submissions, pageSubmissions...)
		}

		if isLastPage(r) {
			break
		}
		page += 1
	}

	return submissions, err
}

func (c *Client) ListMissingSubmissions(opts *ListSubmissionsOpts) ([]*canvasauto.Submission, error) {
	var err error
	missing := make([]*canvasauto.Submission, 0)

	pageInt := 1
	listOpts := &canvasauto.ListMissingSubmissionsParams{
		Page: ptr.To[string]("1"),
	}
	for {
		var r *http.Response
		r, err = c.api.ListMissingSubmissions(c.ctx, opts.UserID, listOpts)
		if err != nil {
			break
		}

		pageMissing := make([]*canvasauto.Submission, 0)
		json.NewDecoder(r.Body).Decode(&pageMissing)

		missing = append(missing, pageMissing...)

		if isLastPage(r) {
			break
		}
		pageInt += 1
		listOpts.Page = ptr.To(strconv.Itoa(pageInt))
	}

	return missing, err
}
