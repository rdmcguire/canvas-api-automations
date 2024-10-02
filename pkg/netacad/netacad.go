package netacad

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Assignment struct {
	Name string
	URL  string
}

// Loads from the parsed file, lazily expecting each line to contain
// a full link (in other words, unformatted crap file)
//
// Adds a prefix to each link, set to "" for no prefix
func LoadAssignmentsHTMLFromFile(file string, prefix string) []Assignment {
	assignmentRegexp := regexp.MustCompile(
		`href="([^"]?launch[^"]+)"[^>]*><div>([^>]+)<`,
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
			URL:  fmt.Sprintf("%s%s", prefix, string(matches[1])),
			Name: strings.Trim(string(matches[2]), " "),
		})
	}

	return assignments
}
