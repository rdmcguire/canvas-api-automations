package netacad

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

type Assignment struct {
	Name string
	URL  string
}

func LoadAssignmentsHTMLFromFile(file string) []Assignment {
	assignmentRegexp := regexp.MustCompile(
		`href="([^"]+)".*instancename"[^>]*>([^>]+)<`,
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
			URL:  string(matches[1]),
			Name: strings.Trim(string(matches[2]), " "),
		})
	}

	return assignments
}
