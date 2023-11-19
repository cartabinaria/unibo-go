// Package degree provides a type to represent a degree course and related
// functions to fetch data from the unibo website.
package degree

import (
	"github.com/csunibo/unibo-go/curriculum"
	"github.com/csunibo/unibo-go/timetable"
)

// Degree represents a degree course.
//
// It can be obtained from UniBo's open data. See the opendata package for more information.
type Degree struct {
	AcademicYear          string // The academic year in which the course is taught
	Code                  int    // The code of the course
	Description           string // The description of the course
	Url                   string // The url of the course website
	Campus                string // The campus where the course is taught
	International         bool   // Whether the course is taught in a language other than Italian
	InternationalTitle    string // The title of the course, if International is true
	InternationalLanguage string // The language in which the course is taught, if International is true
	Fields                string // The fields that the course covers
	Type                  string // The type of the course, e.g. "Laurea triennale"
	DurationInYears       int    // The duration of the course in years, e.g. 3
	OpenForRegistration   string // Whether the course is open for registration
	Languages             string // The languages in which the course is taught
	AccessRequirements    string // The access requirements
	TeachingLocation      string // The main location where the course is taught

	// The id of the course. It is used internally to fetch data from the unibo
	// website. If it is empty, it will be fetched.
	id ID
}

// GetCurricula returns the curricula of the degree for the given year.
// The year must be between 1 and the duration of the degree.
func (d *Degree) GetCurricula(year int) (curriculum.Curricula, error) {
	err := d.fillId()
	if err != nil {
		return nil, err
	}

	curricula, err := curriculum.FetchCurricula(d.id.Type, d.id.Id, year)
	if err != nil {
		return nil, err
	}

	return curricula, nil
}

// GetAllCurricula returns a map of all curricula of the degree.
// The keys are the years of the curricula.
//
// Internally, it calls GetCurricula for each year in a separate goroutine,
// so it is faster than calling GetCurricula for each year.
func (d *Degree) GetAllCurricula() (map[int]curriculum.Curricula, error) {
	err := d.fillId()
	if err != nil {
		return nil, err
	}

	currCh := make(chan curriculum.Curricula)
	errCh := make(chan error)

	for year := 1; year <= d.DurationInYears; year++ {
		go func(year int) {
			curricula, err := curriculum.FetchCurricula(d.id.Type, d.id.Id, year)
			if err != nil {
				errCh <- err
				return
			}
			currCh <- curricula
		}(year)
	}

	curriculaMap := make(map[int]curriculum.Curricula, d.DurationInYears)
	for year := 1; year <= d.DurationInYears; year++ {
		select {
		case curricula := <-currCh:
			curriculaMap[year] = curricula
		case err := <-errCh:
			return nil, err
		}
	}

	return curriculaMap, nil
}

// GetTimetable returns the timetable of the degree for the given year, curriculum and period.
//
// Use GetCurricula or GetAllCurricula to get a curriculum.
// See timetable.FetchTimetable for more information.
func (d *Degree) GetTimetable(
	year int,
	curriculum curriculum.Curriculum,
	period *timetable.Interval,
) (timetable.Timetable, error) {
	err := d.fillId()
	if err != nil {
		return nil, err
	}

	t, err := timetable.FetchTimetable(d.id.Type, d.id.Id, curriculum.Value, year, period)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (d *Degree) fillId() error {
	if d.id != (ID{}) {
		return nil
	}

	id, err := d.ScrapeId()
	if err != nil {
		return err
	}

	d.id = id
	return nil
}
