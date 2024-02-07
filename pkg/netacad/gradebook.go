package netacad

import (
	"encoding/csv"
	"io"
	"os"
	"regexp"
	"strconv"
)

type Gradebook map[Student]*Grades
type Grades map[string]*Grade
type Grade struct {
	Grade      float64
	Percentage float64
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
)

// In Netacad, go to Grades -> Export, select
// Real, Percentage, and Comman delimeter
func LoadGradesFromFile(csvExportFile string) (*Gradebook, error) {
	gradebook := NewGradebook()
	f, err := os.Open(csvExportFile)
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
		gradebook.LoadRow(rowData)

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

func (g *Gradebook) LoadRow(row map[string]string) {
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
