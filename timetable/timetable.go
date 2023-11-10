/*
Package timetable provides methods to retrieve the timetable of a degree.

The timetable is retrieved from the University of Bologna's website. In
particular, every degree has a webpage with the timetable, which is generated
from a JSON endpoint.

For example, the timetable of the degree "Ingegneria Informatica" is available
at https://corsi.unibo.it/laurea/IngegneriaInformatica/orario-lezioni and the
JSON endpoint is available at:

	https://corsi.unibo.it/laurea/IngegneriaInformatica/orario-lezioni/@@orario_reale_json

This package uses the JSON endpoint to retrieve the timetable.
*/
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

// Classroom represents a classroom where an event takes place.
type Classroom struct {
	ResourceDesc string       `json:"des_risorsa"`  // The description of the classroom (e.g. "Classroom 1")
	FloorDesc    string       `json:"des_piano"`    // The floor of the classroom
	BuildingDesc string       `json:"des_edificio"` // The building where the classroom is located
	Raw          RawClassroom `json:"raw"`
}

// RawClassroom represents the raw data of a classroom, as returned by the JSON endpoint.
type RawClassroom struct {
	Enabled      bool     `json:"enabled"`
	Active       bool     `json:"active"`
	Surface      float32  `json:"surface"`
	Blocked      bool     `json:"blocked"`
	EditDate     string   `json:"dataModifica"`
	CreationDate string   `json:"dataCreazione"`
	Seats        int      `json:"numeroPostazioni"`
	Description  string   `json:"descrizione"`
	Building     Building `json:"edificio"`
}

// Building represents a building where a classroom is located.
type Building struct {
	Comune       string `json:"comune"`
	Via          string `json:"via"`
	Provincia    string `json:"provincia"`
	Code         string `json:"codice"`
	CAP          string `json:"cap"`
	Description  string `json:"descrizione"`
	Plesso       string `json:"plesso"`
	Geo          Geo    `json:"geolocalizzazione"`
	CreationDate string `json:"dataCreazione"`
	EditDate     string `json:"dataModifica"`
}

// Geo represents the geographical coordinates of something.
type Geo struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Event represents an event in the timetable.
//
// Usually, an event is a lecture, but it can also be a lab.
type Event struct {
	CodModulo        string       `json:"cod_modulo"`         // The id of the course module
	CalendarInterval string       `json:"periodo_calendario"` // The interval of the event in the calendar
	CodSdoppiamento  string       `json:"cod_sdoppiamento"`   // A code used to identify the lecture, if it is split in multiple parts
	Title            string       `json:"title"`              // The title of the event
	ExtCode          string       `json:"extCode"`            // Unused
	Interval         string       `json:"periodo"`            // The interval of the event
	Teacher          string       `json:"docente"`            // The name of the teacher
	Cfu              int          `json:"cfu"`                // The number of CFU (credits) of the course
	RemoteLearning   bool         `json:"teledidattica"`      // Whether the course is taught remotely
	Teams            string       `json:"teams,omitempty"`    // The link to the Teams meeting. If the course is not taught remotely, this field is omitted
	Start            CalendarTime `json:"start"`              // The start time of the event
	End              CalendarTime `json:"end"`                // The end time of the event
	Classrooms       []Classroom  `json:"aule"`               // The classrooms where the event takes place
}

type Timetable []Event

// Interval represents an interval of time, from Start to End.
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
