package timetable

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	baseUrl       = "https://corsi.unibo.it"
	timetablePath = "/%s/%s/%s/@@orario_reale_json?anno=%d"
)

// Classroom represents a classroom where an event takes place
type Classroom struct {
	Description string `json:"des_risorsa"` // The description of the classroom (e.g. "Classroom 1")
}

type Event struct {
	CodModulo        string       `json:"cod_modulo"` // The id of the course module
	CalendarInterval string       `json:"periodo_calendario"`
	CodSdoppiamento  string       `json:"cod_sdoppiamento"`
	Title            string       `json:"title"`
	ExtCode          string       `json:"extCode"`
	Interval         string       `json:"periodo"`
	Teacher          string       `json:"docente"`         // The name of the teacher
	Cfu              int          `json:"cfu"`             // The number of CFU (credits) of the course
	RemoteLearning   bool         `json:"teledidattica"`   // Whether the course is taught remotely
	Teams            string       `json:"teams,omitempty"` // The link to the Teams meeting. If the course is not taught remotely, this field is omitted
	Start            CalendarTime `json:"start"`           // The start time of the event
	End              CalendarTime `json:"end"`             // The end time of the event
	Classrooms       []Classroom  `json:"aule"`            // The classrooms where the event takes place
}

type Timetable []Event

// Interval represents an interval of time
type Interval struct {
	Start time.Time
	End   time.Time
}

// GetTimetableUrl returns the URL to fetch the timetable for the given course.
//
// See FetchTimetable for the meaning of the parameters.
func GetTimetableUrl(
	courseType, courseId, curriculum string,
	year int,
	interval *Interval,
) string {
	var orarioLang string
	if strings.Contains(courseType, "cycle") {
		orarioLang = "timetable"
	} else {
		orarioLang = "orario-lezioni"
	}

	url := fmt.Sprintf(baseUrl+timetablePath, courseType, courseId, orarioLang, year)

	if curriculum != "" {
		url += fmt.Sprintf("&curricula=%s", curriculum)
	}

	if interval != nil {
		url += fmt.Sprintf("&start=%s", interval.Start.Format("2006-01-02"))
		url += fmt.Sprintf("&end=%s", interval.End.Format("2006-01-02"))
	}

	return url
}

// FetchTimetable retrieves the timetable for the given course.
//
//   - courseType and courseId are the type and id of the course. They must not be empty.
//   - curriculum is the curriculum of the course. It can be empty.
//   - year is the academic year of the timetable.
//   - interval is the wanted interval of the timetable. It can be nil.
func FetchTimetable(
	courseType, courseId, curriculum string,
	year int,
	interval *Interval,
) (timetable Timetable, err error) {
	url := GetTimetableUrl(courseType, courseId, curriculum, year, interval)

	res, err := http.Get(url)
	if err != nil {
		return
	}

	err = json.NewDecoder(res.Body).Decode(&timetable)
	if err != nil {
		return
	}

	err = res.Body.Close()
	if err != nil {
		return
	}

	return
}
