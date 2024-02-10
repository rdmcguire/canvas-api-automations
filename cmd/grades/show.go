package grades

import (
	"sync"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var gradesShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"list", "ls", "s"},
	Short:   "Shows graded submissions given any provided flags",
	Long: `Run the show command with --email, --moduleID, or --assignmentID flag(s)
These can each be provided multiple times, or comma-delimeted`,
	Run: execShowGradesCmd,
}

// Cache users to prevent repeat lookups
var (
	gradedOnly, ungradedOnly bool
	wg                       sync.WaitGroup
)

func execShowGradesCmd(cmd *cobra.Command, args []string) {
	log := util.Logger(cmd)

	moduleIDs, _ := cmd.Flags().GetIntSlice("moduleID")
	assignmentIDs, _ := cmd.Flags().GetIntSlice("assignmentID")
	gradedOnly, _ = cmd.Flags().GetBool("gradedOnly")
	ungradedOnly, _ = cmd.Flags().GetBool("ungradedOnly")

	log.Debug().
		Ints("moduleIDs", moduleIDs).
		Ints("assignmentIDs", assignmentIDs).
		Bool("gradedOnly", gradedOnly).
		Bool("ungradedOnly", ungradedOnly).
		Str("courseID", util.GetCourseIdStr(cmd)).
		Msg("Listing grades")

	assignments := listAssignments(cmd)

	// Start smashing against assignments in goroutines
	for _, a := range assignments {
		wg.Add(1)
		go showAssignment(cmd, a)
	}
	wg.Wait()
}

func showAssignment(cmd *cobra.Command, assignment *canvasauto.Assignment) {
	defer wg.Done()
	log := util.Logger(cmd)
	log.Debug().
		Str("assignment", canvas.StrOrNil(assignment.Name)).
		Int("ID", *assignment.Id).
		Msg("Found assignment")
	for _, s := range listSubmissions(cmd, assignment) {
		if gradedOnly && *s.WorkflowState != "graded" {
			continue
		} else if ungradedOnly && *s.WorkflowState == "graded" {
			continue
		}
		user := util.Client(cmd).GetUserById(util.GetCourseIdStr(cmd), *s.UserId)
		if user == nil {
			panic(*s)
		}
		log.Info().
			Str("assignment", *assignment.Name).
			Any("score", s.Score).
			Str("state", canvas.StrOrNil(s.WorkflowState)).
			Str("user", *user.Email).
			Msg("Found Grade")
	}
}

func listSubmissions(cmd *cobra.Command, assignment *canvasauto.Assignment) []*canvasauto.Submission {
	client := util.Client(cmd)
	log := util.Logger(cmd)
	emails, _ := cmd.Flags().GetStringArray("email")

	submissions, err := client.ListAssignmentSubmissions(&canvas.ListSubmissionsOpts{
		CourseID:   util.GetCourseIdStr(cmd),
		Assignment: assignment,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed listing submissions")
	}

	filteredSubmissions := make([]*canvasauto.Submission, 0, len(submissions))
	for _, s := range submissions {
		user := client.GetUserById(util.GetCourseIdStr(cmd), *s.UserId)
		if user == nil {
			log.Debug().
				Int("userID", *s.UserId).
				Any("submission", *s).
				Msg("Failed to find user from assignment, ignoring submission")
			continue
		}

		if len(emails) > 0 && !slices.Contains(emails, *user.Email) {
			continue
		}
		filteredSubmissions = append(filteredSubmissions, s)
	}

	return filteredSubmissions
}

func listAssignments(cmd *cobra.Command) []*canvasauto.Assignment {
	moduleIDs, _ := cmd.Flags().GetIntSlice("moduleID")
	assignmentIDs, _ := cmd.Flags().GetIntSlice("assignmentID")
	assignments := make([]*canvasauto.Assignment, 0)

	assignmentsByModule := util.Client(cmd).ListAssignmentsByModule(
		util.GetCourseIdStr(cmd),
		moduleIDs...,
	)

	// Filter if assignments provided
	for _, module := range assignmentsByModule {
		for _, assignment := range module {
			if len(assignmentIDs) > 0 && !slices.Contains(assignmentIDs, *assignment.Id) {
				continue
			}
			assignments = append(assignments, assignment)
		}
	}

	return assignments
}

func init() {
	gradesShowCmd.Flags().IntSliceP("moduleID", "m", []int{}, "Filter by module ID")
	gradesShowCmd.Flags().IntSliceP("assignmentID", "a", []int{}, "Filter by assignment ID")
	gradesShowCmd.Flags().BoolP("gradedOnly", "g", false, "Restrict to graded grades only")
	gradesShowCmd.Flags().BoolP("ungradedOnly", "G", false, "Restrict to ungraded grades only")
}
