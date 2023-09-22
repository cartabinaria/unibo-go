package degree

import (
	"fmt"

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
}

// GetCurricula returns the curricula of the degree for the given year.
// The year must be between 1 and the duration of the degree.
func (d Degree) GetCurricula(year int) (curriculum.Curricula, error) {
	id, err := d.ScrapeId()
	if err != nil {
		return nil, err
	}

	curricula, err := curriculum.FetchCurricula(id.Type, id.Id, year)
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
func (d Degree) GetAllCurricula() (map[int]curriculum.Curricula, error) {
	id, err := d.ScrapeId()
	if err != nil {
		return nil, fmt.Errorf("could not get course website id: %w", err)
	}

	currCh := make(chan curriculum.Curricula)
	errCh := make(chan error)

	for year := 1; year <= d.DurationInYears; year++ {
		go func(year int) {
			curricula, err := curriculum.FetchCurricula(id.Type, id.Id, year)
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

func (d Degree) GetTimetable(year int, curriculum curriculum.Curriculum, period *timetable.Interval) (timetable.Timetable, error) {
	id, err := d.ScrapeId()
	if err != nil {
		return nil, err
	}

	t, err := timetable.FetchTimetable(id.Type, id.Id, curriculum.Value, year, period)
	if err != nil {
		return nil, err
	}

	return t, nil
}
