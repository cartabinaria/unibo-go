// Package curriculum contains functions to fetch curricula from the unibo
// website
package curriculum

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	baseCurriculaIt = "https://corsi.unibo.it/%s/%s/orario-lezioni/@@available_curricula?anno=%d&curricula="
	baseCurriculaEn = "https://corsi.unibo.it/%s/%s/timetable/@@available_curricula?anno=%d&curricula="
)

type (
	Curriculum struct {
		Selected bool   `json:"selected"`
		Value    string `json:"value"`
		Label    string `json:"label"`
	}
	Curricula []Curriculum
)

func GetCurriculaUrl(courseType, courseId string, year int) string {
	if strings.Contains(courseType, "cycle") {
		return fmt.Sprintf(baseCurriculaEn, courseType, courseId, year)
	} else {
		return fmt.Sprintf(baseCurriculaIt, courseType, courseId, year)
	}
}

func FetchCurricula(courseType, courseId string, year int) (curricula Curricula, err error) {
	url := GetCurriculaUrl(courseType, courseId, year)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&curricula)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return
}
