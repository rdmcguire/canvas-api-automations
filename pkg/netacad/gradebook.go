package netacad

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
)

const (
	colName  = "NAME"
	colEmail = "EMAIL"
)

var (
	gradeRegexp          = regexp.MustCompile(`(.*) \((Real|Percentage)\)`)
	pcntGradeRegexp      = regexp.MustCompile(`([0-9.]+) ?%`)
	isGradeRegexp        = regexp.MustCompile(`^[0-9]+(\.[0-9]+)? ?%?`)
	isTotalRegexp        = regexp.MustCompile(`total$`)
	pointsPossibleRegexp = regexp.MustCompile(`^Points? Possible$`)

	pointsPossible map[string]float64
)

type (
	Gradebook map[Student]*Grades
	Grades    map[string]*Grade
	Grade     struct {
		Assignment  *canvasauto.Assignment
		Module      *canvasauto.Module
		User        *canvasauto.User
		Submissions []*canvasauto.Submission
		Grade       float64
		Percentage  float64
	}
)

type Student struct {
	First string
	Last  string
	Email string
	ID    int64
}

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
		return assignments, errors.New("no headers found in grade export")
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
		} else if pointsPossible == nil {
			pointsPossible, err = loadPointsPossible(row, headers)
			if err != nil {
				return nil, err
			}
		}

		rowData := rowToMap(headers, row)

		// Filter unwanted emails
		if len(opts.Emails) > 0 {
			if !hasEmail(opts.Emails, rowData[colEmail]) {
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

func loadPointsPossible(row []string, headers []string) (map[string]float64, error) {
	pointsPossible = make(map[string]float64, len(headers)-2)
	var err error
	if !pointsPossibleRegexp.Match([]byte(row[0])) {
		return pointsPossible,
			fmt.Errorf("unexpected field %s for points possible stupidity", row[0])
	}

	for field, points := range row {
		if field < 2 {
			continue
		} else if points == "" {
			continue
		}

		possible, parseErr := strconv.ParseFloat(points, 64)
		if parseErr != nil {
			return nil, parseErr
		}
		pointsPossible[headers[field]] = possible
	}
	return pointsPossible, err
}

func hasEmail(filter []string, email string) bool {
	for _, eml := range filter {
		if strings.EqualFold(eml, email) {
			return true
		}
	}

	return false
}

func rowToMap(headers []string, row []string) map[string]string {
	data := make(map[string]string, len(row))
	for i, v := range row {
		data[headers[i]] = v
	}
	return data
}

func (g *Gradebook) loadRow(row map[string]string) {
	if row[colName] == "" {
		return
	}

	var first, last string
	name := strings.Split(row[colName], " ")
	if len(name) > 0 {
		first = strings.Join(name[:len(name)-1], " ")
		last = name[len(name)-1]
	}
	student := Student{
		First: first,
		Last:  last,
		Email: row[colEmail],
	}

	var grades *Grades

	if (*g)[student] == nil {
		grades = NewGrades()
		(*g)[student] = grades
	}

	for key, grade := range row {
		if !isGradeRegexp.MatchString(grade) {
			continue
		}
		possible, ok := pointsPossible[key]
		if !ok {
			log.Warn().Str("key", key).Str("grade", grade).Msg("encountered unknown assignment")
			continue
		}
		gradeFloat, err := strconv.ParseFloat(grade, 64)
		if err != nil {
			log.Err(err).Str("key", key).Str("grade", grade).Msg("unable to parse grade")
			continue
		}
		pcnt := gradeFloat / possible
		grades.RecordGrade(key, gradeFloat, pcnt)
	}
}

func (g *Grades) RecordGrade(item string, grade float64, pcnt float64) {
	if (*g)[item] == nil {
		(*g)[item] = &Grade{}
	}
	(*g)[item].Grade = grade
	(*g)[item].Percentage = pcnt
}

// NOTE: This is legacy code, see the comment in GradeItemFromHeader
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
// NOTE: this used to be required in `old netacad` but new one is
// even more dumb than the old one. Keeping in-case the Netacad clowns
// decide to keep making stupid changes that add no value just to annoy
// the piss out of anyone trying to automate with their obnoxious
// data and lack of API
func GradeItemFromHeader(header string) (string, string) {
	if parts := gradeRegexp.FindStringSubmatch(header); len(parts) == 3 {
		return parts[1], parts[2]
	}
	return "", ""
}

func NewGrades() *Grades {
	var grades Grades = make(map[string]*Grade, 0)
	return &grades
}

func NewGradebook() *Gradebook {
	var gradebook Gradebook = make(map[Student]*Grades, 0)
	return &gradebook
}
