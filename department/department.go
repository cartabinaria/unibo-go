package department

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

const (
	maxTeachers    = "2000"
	rootUrl        = "https://%s.unibo.it/it/dipartimento/persone/docenti-e-ricercatori"
	departmentsUrl = "https://www.unibo.it/it/ateneo/sedi-e-strutture/dipartimenti"
)

var (
	departmentsRegex = regexp.MustCompile("<a class=\"internal-link\" href=\"https://(.+).unibo.it/it\".+>(.+)</a>")
)

// getDepartmentTeacherUrl returns the URL to fetch the list of teachers for the given department.
func getDepartmentTeacherUrl(department string) string {
	return fmt.Sprintf(rootUrl, department) + "?pagesize=" + maxTeachers
}

// Department represents a department of the university.
type Department struct {
	Url  string // The url is the url of the department website, e.g. "https://disi.unibo.it/it"
	Name string // The name is the name of the department, e.g. "Informatica - Scienza e Ingegneria"
	Code string // The code is the subdomain of the department website, e.g. "disi" for "https://disi.unibo.it/it"
}

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

		dep := Department{Code: code, Name: name, Url: "https://" + code + ".unibo.it/it"}
		deps = append(deps, dep)
	}

	return deps, nil
}
