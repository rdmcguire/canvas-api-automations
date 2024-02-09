package students

import (
	"fmt"
	"strings"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd/util"
	"github.com/spf13/cobra"
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

func execStudentsExportCmd(cmd *cobra.Command, args []string) {
	client := util.Client(cmd)
	log := util.Logger(cmd)

	users := client.ListUsersInCourse(util.GetCourseIdStr(cmd), "")
	log.Debug().Int("numStudents", len(users)).
		Msg("Successfully retrieved students from course")

	fmt.Println("First Name,Last Name,Email Address,Student ID")
	for _, user := range users {
		nameParts := strings.Split(*user.SortableName, ", ")
		if len(nameParts) < 2 {
			continue
		}
		fmt.Printf("%s,%s,%s,%d\n", nameParts[1], nameParts[0], *user.Email, user.Id)
	}
}
