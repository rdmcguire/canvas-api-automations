package grades

import (
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/students"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/netacad"
	"github.com/spf13/cobra"
)

var GradesCmd = &cobra.Command{
	Use:              "grades",
	Aliases:          []string{"g"},
	Short:            "sub-commands for grading",
	PersistentPreRun: util.EnsureCourseId,
}

// Used to cache assignment data
var assignmentCache *util.AssignmentCache

func init() {
	GradesCmd.AddCommand(gradesDumpCmd)
	GradesCmd.AddCommand(gradesPushCmd)
	GradesCmd.AddCommand(gradesShowCmd)
	GradesCmd.AddCommand(gradesBulkCmd)
	GradesCmd.AddCommand(gradesFindCmd)

	GradesCmd.PersistentFlags().StringArray("email", []string{}, "Restrict to a provided email addresses")
	GradesCmd.RegisterFlagCompletionFunc("email", students.ValidateEmailArg)
}

func mustLoadGrades(cmd *cobra.Command, file string) *netacad.Gradebook {
	log := util.Logger(cmd)

	loadOpts := &netacad.LoadGradesFromFileOpts{
		File:           file,
		WithGradesOnly: true,
	}
	if emails, _ := cmd.Flags().GetStringArray("email"); len(emails) > 0 {
		loadOpts.Emails = emails
	}

	gradebook, err := netacad.LoadGradesFromFile(loadOpts)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load netacad grade export csv")
	}

	return gradebook
}
