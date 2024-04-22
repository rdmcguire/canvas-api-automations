package netacad

import (
	"bufio"
	"os"
	"regexp"
)

type Assignment struct {
	Name string
	Url  string
}

func LoadAssignmentsHtmlFromFile(file string) []Assignment {
	assignmentRegexp := regexp.MustCompile(
		`class="activity.*href="([^"]+)".*instancename"[^>]*>([^>]+)<`,
	)

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	assignments := make([]Assignment, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		matches := assignmentRegexp.FindSubmatch(scanner.Bytes())
		if len(matches) < 3 {
			continue
		}
		assignments = append(assignments, Assignment{
			Url:  string(matches[1]),
			Name: string(matches[2]),
		})
	}

	return assignments
}
