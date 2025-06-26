// SPDX-FileCopyrightText: 2024 - 2025 Eyad Issa <eyadlorenzo@gmail.com>
// SPDX-FileCopyrightText: 2025 Samuele Musiani <samu@teapot.ovh>
//
// SPDX-License-Identifier: MIT

package exams

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
)

var (
	baseUrl = "https://corsi.unibo.it/%s/%s/appelli"
)

type Exam struct {
	SubjectCode   string
	SubjectName   string
	Teacher       string
	Date          time.Time
	Type          string
	Location      string
	Subscriptions string
}

var duplicatedSpaceRemover = regexp.MustCompile(`\s+`)

func GetExams(courseType, courseId string) ([]Exam, error) {
	return GetExamsForSubject(courseType, courseId, "")
}

// subjectsPerPage is the number of subjects that are shown per page on the website
const subjectsPerPage = 20

func GetExamsForSubject(courseType, courseId, subjectName string) ([]Exam, error) {
	var exams []Exam

	start := 0
	for {
		url := fmt.Sprintf(baseUrl, courseType, courseId)

		// if we are looking for a specific subject, we need to specify the appelli parameter
		if subjectName != "" {
			url = fmt.Sprintf("%s?appelli=%s", url, subjectName)
		}

		// if we are not at the first page, we need to specify the start parameter
		if start > 1 {
			url = fmt.Sprintf("%s?b_start:int=%d", url, start)
		}

		document, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch exams from url %s: %w", url, err)
		}

		defer func() {
			if err := document.Body.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "unable to close response body: %v", err)
			}
		}()

		e, err := parseExamsHtml(document.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to parse exams from url %s: %w", url, err)
		}

		if len(e) == 0 { // no more exams
			break
		}
		start += subjectsPerPage
		exams = append(exams, e...)
	}

	return exams, nil
}

func parseExamsHtml(r io.Reader) ([]Exam, error) {
	node, err := htmlquery.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("unable to load url: %w", err)
	}

	tabs := htmlquery.Find(node, "//h3[@role='tab']")
	panels := htmlquery.Find(node, "//div[@role='tabpanel']")

	if len(tabs) != len(panels) {
		return nil, fmt.Errorf("unable to find the same number of tab panels and tabs while parsing exams. maybe the html structure has changed")
	}

	exams := make([]Exam, 0)

	for i := 0; i < len(tabs); i++ {

		panelNode := tabs[i]

		subjNode := htmlquery.FindOne(panelNode, "/a")
		if subjNode == nil {
			return nil, fmt.Errorf("unable to find subject while parsing exams. maybe the html structure has changed")
		}

		subjCodeNode := htmlquery.FindOne(panelNode, "//span[@class='code']")
		if subjCodeNode == nil {
			return nil, fmt.Errorf("unable to find code while parsing exams. maybe the html structure has changed")
		}

		subjCodeStr := htmlquery.InnerText(subjCodeNode)

		subjTeacherNode := htmlquery.FindOne(subjNode, "//span[@class='docente']")
		if subjTeacherNode == nil {
			return nil, fmt.Errorf("unable to find teacher while parsing exams. maybe the html structure has changed")
		}
		subjTeacher := htmlquery.InnerText(subjTeacherNode)
		subjTeacher = strings.TrimSpace(subjTeacher)

		spaceRemover := strings.NewReplacer("\n", "", "\t", "", subjCodeStr, "", subjTeacher, "")

		title := htmlquery.InnerText(subjNode)
		title = spaceRemover.Replace(title)
		title = duplicatedSpaceRemover.ReplaceAllString(title, " ")
		title = strings.TrimSpace(title)

		tabNode := panels[i]

		examsNodes := htmlquery.Find(tabNode, "/table")
		if len(examsNodes) == 0 {
			return nil, fmt.Errorf("unable to find exams while parsing exams. maybe the html structure has changed")
		}

		cleanText := func(text string) string {
			text = spaceRemover.Replace(text)
			text = duplicatedSpaceRemover.ReplaceAllString(text, " ")
			text = strings.TrimSpace(text)
			return text
		}

		for _, examNode := range examsNodes {

			dateNode := htmlquery.FindOne(examNode, "//tr[1]/td[1]")
			if dateNode == nil {
				return nil, fmt.Errorf("unable to find date while parsing exams. maybe the html structure has changed")
			}
			date := htmlquery.InnerText(dateNode)
			date = cleanText(date)

			listaIscrizioniNode := htmlquery.FindOne(examNode, "//tr[2]/td[1]")
			if listaIscrizioniNode == nil {
				return nil, fmt.Errorf("unable to find lista iscrizioni while parsing exams. maybe the html structure has changed")
			}
			listaIscrizioni := htmlquery.InnerText(listaIscrizioniNode)
			listaIscrizioni = cleanText(listaIscrizioni)

			tipoProvaNode := htmlquery.FindOne(examNode, "//tr[3]/td[1]")
			if tipoProvaNode == nil {
				return nil, fmt.Errorf("unable to find tipo prova while parsing exams. maybe the html structure has changed")
			}
			tipoProva := htmlquery.InnerText(tipoProvaNode)
			tipoProva = cleanText(tipoProva)

			luogoNode := htmlquery.FindOne(examNode, "//tr[4]/td[1]")
			if luogoNode == nil {
				return nil, fmt.Errorf("unable to find luogo while parsing exams. maybe the html structure has changed")
			}
			luogo := htmlquery.InnerText(luogoNode)
			luogo = cleanText(luogo)

			parsedDate, err := ParseItalianDate(date)
			if err != nil {
				return nil, fmt.Errorf("unable to parse date '%s': %w", date, err)
			}

			exam := Exam{
				SubjectCode:   subjCodeStr,
				SubjectName:   title,
				Teacher:       subjTeacher,
				Date:          parsedDate,
				Type:          tipoProva,
				Location:      luogo,
				Subscriptions: listaIscrizioni,
			}

			exams = append(exams, exam)
		}
	}

	return exams, nil
}

// ParseItalianDate parses a date string in the format "mese anno ore hh:mm" (e.g., "febbraio 2025 ore 09:30")
// and returns a time.Time object.
func ParseItalianDate(input string) (time.Time, error) {
	// Map Italian months to their English counterparts
	months := map[string]string{
		"gennaio":   "January",
		"febbraio":  "February",
		"marzo":     "March",
		"aprile":    "April",
		"maggio":    "May",
		"giugno":    "June",
		"luglio":    "July",
		"agosto":    "August",
		"settembre": "September",
		"ottobre":   "October",
		"novembre":  "November",
		"dicembre":  "December",
	}

	// Replace Italian month with English month
	for italian, english := range months {
		if strings.Contains(input, italian) {
			input = strings.Replace(input, italian, english, 1)
			break
		}
	}

	// Replace "ore" with space for compatibility
	input = strings.Replace(input, "ore ", "", 1)

	timezone, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		return time.Time{}, fmt.Errorf("could not load italian timezone: %w", err)
	}

	// Parse the time
	parsedTime, err := time.ParseInLocation("02 January 2006 15:04", strings.TrimSpace(input), timezone)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}
