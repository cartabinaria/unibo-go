package department

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

// Department represents a department of the university.
type Department struct {
	Name string // The name is the name of the department, e.g. "Informatica - Scienza e Ingegneria"
	Code string // The code is the subdomain of the department website, e.g. "disi" for "https://disi.unibo.it/it"
}

func (d Department) Url() string                       { return "https://" + d.Code + ".unibo.it/it" }
func (d Department) FetchTeachers() ([]Teacher, error) { return FetchTeachers(d.Code) }
func (d Department) GetTeachersUrl() string            { return getDepartmentTeacherUrl(d.Code) }

const departmentsUrl = "https://www.unibo.it/it/ateneo/sedi-e-strutture/dipartimenti"

var departmentsRegex = regexp.MustCompile("<a class=\"internal-link\" href=\"https://(.+).unibo.it/it\".+>(.+)</a>")

// FetchDepartments retrieves the list of departments of the university.
//
// It gets the list from the university website via HTTP and then applies a regex
// to parse the HTML.
func FetchDepartments() ([]Department, error) {
	res, err := http.Get(departmentsUrl)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	matches := departmentsRegex.FindAllSubmatch(body, -1)
	deps := make([]Department, 0, len(matches))
	for _, match := range matches {
		code := string(match[1])
		name := string(match[2])

		dep := Department{
			Code: code,
			Name: name,
		}
		deps = append(deps, dep)
	}

	return deps, nil
}

const (
	maxTeachers = "2000"
	teachersUrl = "https://%s.unibo.it/it/dipartimento/persone/docenti-e-ricercatori"
)

// getDepartmentTeacherUrl returns the URL to fetch the list of teachers for the given department.
func getDepartmentTeacherUrl(department string) string {
	return fmt.Sprintf(teachersUrl, department) + "?pagesize=" + maxTeachers
}
