package department

import (
	"io"
	"net/http"
	"regexp"
)

// Teacher represents a teacher.
type Teacher struct {
	Username string
}

// GetWebsite returns the website of the teacher.
func (t Teacher) GetWebsite() string {
	return "https://www.unibo.it/sitoweb/" + t.Username
}

var teachersRegex = regexp.MustCompile("<a href=\"https://www.unibo.it/sitoweb/([a-z0-9.]+)\"")

// FetchTeachers retrieves the list of teachers for the given department.
func FetchTeachers(departmentCode string) ([]Teacher, error) {
	url := getDepartmentTeacherUrl(departmentCode)

	res, err := http.Get(url)
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

	var teachers []Teacher
	matches := teachersRegex.FindAllSubmatch(body, -1)
	for _, match := range matches {
		teacher := Teacher{Username: string(match[1])}
		teachers = append(teachers, teacher)
	}

	return teachers, nil
}
