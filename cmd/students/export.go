package students

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
)

var studentsExportCmd = &cobra.Command{
	Use:     "export",
	Aliases: []string{"dump", "csv"},
	Args:    cobra.NoArgs,
	Run:     execStudentsExportCmd,
	Short:   "Exports Canvas Course Students",
	Long: `This command will write a list of students in
the csv format expected by Netacad. Redirect to your file
( > somefile.csv) and load directly into your new Netacad
course using their bulk import.`,
}

const (
	fmtQuotedWithID    = `"%s","%s","%s","%d"`
	fmtWithID          = `%s,%s,%s,%d`
	fmtQuotedWithoutID = `"%s","%s","%s"`
	fmtWithoutID       = `%s,%s,%s`

	flagID             = "withID"
	flagQuote          = "withQuotes"
	flagStripFirstName = "stripFirst"
)

var (
	withID             bool
	withQuotes         bool
	withFirstNameStrip bool

	regexReplaceWithSpace = regexp.MustCompile(`[_|]`)
)

func execStudentsExportCmd(cmd *cobra.Command, args []string) {
	client := util.Client(cmd)
	log := util.Logger(cmd)

	withID, _ = cmd.Flags().GetBool(flagID)
	withQuotes, _ = cmd.Flags().GetBool(flagQuote)
	withFirstNameStrip, _ = cmd.Flags().GetBool(flagStripFirstName)

	users := client.ListUsersInCourse(util.GetCourseIdStr(cmd), "")
	log.Debug().Int("numStudents", len(users)).
		Msg("Successfully retrieved students from course")

	fmt.Println(headerString())

	for _, user := range users {
		if userStr := userString(user); userStr != nil {
			fmt.Println(*userStr)
		} else {
			log.Fatal().Any("user", *user).Msg("unable to generate csv output for user")
		}
	}
}

func userString(user *canvasauto.User) *string {
	nameParts := strings.Split(*user.SortableName, ", ")
	if len(nameParts) < 2 {
		return nil
	}

	email := strings.ReplaceAll(*user.Email, "\"", "'")
	first := strings.ReplaceAll(nameParts[1], "\"", "'")
	last := strings.ReplaceAll(nameParts[0], "\"", "'")

	first = regexReplaceWithSpace.ReplaceAllLiteralString(first, " ")
	last = regexReplaceWithSpace.ReplaceAllLiteralString(last, " ")

	if withFirstNameStrip {
		first = strings.Split(first, " ")[0]
	}

	var userStr string

	if withQuotes {
		if withID {
			userStr = fmt.Sprintf(fmtQuotedWithID, first, last, email, user.Id)
		} else {
			userStr = fmt.Sprintf(fmtQuotedWithoutID, first, last, email)
		}
	} else {
		if withID {
			userStr = fmt.Sprintf(fmtWithID, first, last, email, user.Id)
		} else {
			userStr = fmt.Sprintf(fmtWithoutID, first, last, email)
		}
	}

	return &userStr
}

func headerString() string {
	if withQuotes {
		if withID {
			return "\"First Name\",\"Last Name\",\"Email Address\",\"Student ID\""
		} else {
			return "\"First Name\",\"Last Name\",\"Email Address\""
		}
	} else {
		if withID {
			return "First Name,Last Name,Email Address,Student ID"
		} else {
			return "First Name,Last Name,Email Address"
		}
	}
}

func init() {
	studentsExportCmd.Flags().Bool(flagID, false, "Include Student ID in export")
	studentsExportCmd.Flags().Bool(flagQuote, true, "Quote all fields")
	studentsExportCmd.Flags().Bool(flagStripFirstName, false,
		"Use only the first part of first name, strip anything after space")
}
