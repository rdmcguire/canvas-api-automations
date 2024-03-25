package netacad

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"golang.org/x/exp/slices"
)

type Gradebook map[Student]*Grades
type Grades map[string]*Grade
type Grade struct {
	Grade      float64
	Percentage float64
	// These fields should be updated once
	// matching information is found in the canvas API
	Assignment  *canvasauto.Assignment   // Assignment (to be found)
	Submissions []*canvasauto.Submission // Submissions (to be loaded)
	Module      *canvasauto.Module       // Module (to be found)
	User        *canvasauto.User         // User (to be found)
}

type Student struct {
	ID    int64
	First string
	Last  string
	Email string
}

var (
	gradeRegexp     = regexp.MustCompile(`(.*) \((Real|Percentage)\)`)
	pcntGradeRegexp = regexp.MustCompile(`([0-9.]+) ?%`)
	isGradeRegexp   = regexp.MustCompile(`^[0-9]+(\.[0-9]+)? ?%?`)
	isTotalRegexp   = regexp.MustCompile(`total$`)
)

type LoadGradesFromFileOpts struct {
	File           string   // Path to Netacad csv export
	Emails         []string // Optional email filter
	WithGradesOnly bool     // Only return students with gradeable items
}

// Returns a list of assignments from the given file
// Filters out totals columns
func Assignments(opts *LoadGradesFromFileOpts) ([]string, error) {
	assignments := make([]string, 0)

	f, err := os.Open(opts.File)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := csv.NewReader(f)
	headers, err := data.Read() // Read one row
	if err != nil {
		return assignments, err
	} else if len(headers) < 1 {
		return assignments, errors.New("No headers found in grade export")
	}

	for _, h := range headers {
		name, _ := GradeItemFromHeader(h)
		if name == "" {
			continue
		} else if isTotalRegexp.Match([]byte(name)) {
			continue
		} else if !slices.Contains(assignments, name) {
			assignments = append(assignments, name)
		}
	}

	return assignments, nil
}

// In Netacad, go to Grades -> Export, select
// Real, Percentage, and Comman delimeter
func LoadGradesFromFile(opts *LoadGradesFromFileOpts) (*Gradebook, error) {
	gradebook := NewGradebook()
	f, err := os.Open(opts.File)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := csv.NewReader(f)

	headers := make([]string, 0)
	for {
		row, err := data.Read()
		if len(headers) == 0 {
			headers = row
			continue
		}

		rowData := rowToMap(headers, row)

		// Filter unwanted emails
		if len(opts.Emails) > 0 {
			if !slices.Contains(opts.Emails, rowData["Email address"]) {
				goto Next
			}
		}
		gradebook.loadRow(rowData)

	Next:
		if err == io.EOF {
			break
		} else if err != nil {
			return gradebook, err
		}
	}

	return gradebook, nil
}

func rowToMap(headers []string, row []string) map[string]string {
	data := make(map[string]string, len(row))
	for i, v := range row {
		data[headers[i]] = v
	}
	return data
}

func (g *Gradebook) loadRow(row map[string]string) {
	id, err := strconv.ParseInt(row["ID number"], 10, 64)
	if err != nil {
		return
	}

	student := Student{
		ID:    id,
		First: row["First name"],
		Last:  row["Surname"],
		Email: row["Email address"],
	}

	if (*g)[student] == nil {
		(*g)[student] = NewGrades()
	}

	for key, grade := range row {
		item, itemType := GradeItemFromHeader(key)
		if item == "" || itemType == "" {
			continue
		}
		(*g)[student].Record(item, itemType, grade)
	}
}

func (g *Grades) Record(item string, itemType string, grade string) {
	if !isGradeRegexp.MatchString(grade) {
		return
	}

	if (*g)[item] == nil {
		(*g)[item] = &Grade{}
	}

	switch itemType {
	case "Real":
		grade, err := strconv.ParseFloat(grade, 64)
		if err != nil {
			panic(err)
		}
		(*g)[item].Grade = grade

	case "Percentage":
		if parts := pcntGradeRegexp.FindStringSubmatch(grade); len(parts) == 2 {
			grade, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				panic(err)
			}
			(*g)[item].Percentage = grade
		}
	}
}

func (g *Grades) Count() int {
	return len(*g)
}

// Returns gradeable column from header string.
// Returns name of the item and type (Real|Percentage) separately
func GradeItemFromHeader(header string) (string, string) {
	if parts := gradeRegexp.FindStringSubmatch(header); len(parts) == 3 {
		return parts[1], parts[2]
	}
	return "", ""
}

func NewGrades() *Grades {
	var grades Grades
	grades = make(map[string]*Grade, 0)
	return &grades
}

func NewGradebook() *Gradebook {
	var gradebook Gradebook
	gradebook = make(map[Student]*Grades, 0)
	return &gradebook
}
