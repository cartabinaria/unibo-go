// SPDX-FileCopyrightText: 2023 Eyad Issa <eyadlorenzo@gmail.com>
// SPDX-FileCopyrightText: 2024 Samuele Musiani <samu@teapot.ovh>

// Package curriculum contains functions to fetch curricula from the unibo
// website
package curriculum

import (
	"encoding/json"
	"fmt"
	"io"
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if strings.Contains(string(body), "error") {
		return nil, fmt.Errorf("Unibo website returned an error for url: %s", url)
	}

	err = json.Unmarshal(body, &curricula)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return
}
