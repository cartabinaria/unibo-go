package exams

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
)

var (
	baseUrl = "https://corsi.unibo.it/%s/%s/appelli"
)

type Exam struct {
	SubjectCode   int
	SubjectName   string
	Teacher       string
	Date          string
	Type          string
	Location      string
	Subscriptions string
}

var duplicatedSpaceRemover = regexp.MustCompile(`\s+`)

func GetExams(courseType, courseId string) ([]Exam, error) {
	return GetExamsForSubject(courseType, courseId, "")
}

func GetExamsForSubject(courseType, courseId, subjectName string) ([]Exam, error) {
	url := fmt.Sprintf(baseUrl, courseType, courseId)
	if subjectName != "" {
		url = fmt.Sprintf("%s?appelli=%s", url, subjectName)
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

	node, err := htmlquery.Parse(document.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to load url: %w", err)
	}

	subjNodes := htmlquery.Find(node, "//div[@role='tablist']")

	exams := make([]Exam, 0)

	for _, appelloNode := range subjNodes {

		subjNode := htmlquery.FindOne(appelloNode, "//h3/a[@class='openclose-appelli']")
		if subjNode == nil {
			return nil, fmt.Errorf("unable to find subject while parsing exams. maybe the html structure has changed")
		}

		subjCodeNode := htmlquery.FindOne(subjNode, "//span[@class='code']")
		if subjCodeNode == nil {
			return nil, fmt.Errorf("unable to find code while parsing exams. maybe the html structure has changed")
		}

		subjCodeStr := htmlquery.InnerText(subjCodeNode)
		subjCode, err := strconv.Atoi(subjCodeStr)
		if err != nil {
			return nil, fmt.Errorf("unable to convert code to int: %w", err)
		}

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

		examsNodes := htmlquery.Find(appelloNode, "//table[@class='single-item']")
		if len(examsNodes) == 0 {
			return nil, fmt.Errorf("unable to find exams while parsing exams. maybe the html structure has changed")
		}

		for _, examNode := range examsNodes {

			dateNode := htmlquery.FindOne(examNode, "//tr[1]/td[1]")
			if dateNode == nil {
				return nil, fmt.Errorf("unable to find date while parsing exams. maybe the html structure has changed")
			}
			date := htmlquery.InnerText(dateNode)
			date = strings.TrimSpace(date)

			listaIscrizioniNode := htmlquery.FindOne(examNode, "//tr[2]/td[1]")
			if listaIscrizioniNode == nil {
				return nil, fmt.Errorf("unable to find lista iscrizioni while parsing exams. maybe the html structure has changed")
			}
			listaIscrizioni := htmlquery.InnerText(listaIscrizioniNode)
			listaIscrizioni = strings.TrimSpace(listaIscrizioni)

			tipoProvaNode := htmlquery.FindOne(examNode, "//tr[3]/td[1]")
			if tipoProvaNode == nil {
				return nil, fmt.Errorf("unable to find tipo prova while parsing exams. maybe the html structure has changed")
			}
			tipoProva := htmlquery.InnerText(tipoProvaNode)
			tipoProva = strings.TrimSpace(tipoProva)

			luogoNode := htmlquery.FindOne(examNode, "//tr[4]/td[1]")
			if luogoNode == nil {
				return nil, fmt.Errorf("unable to find luogo while parsing exams. maybe the html structure has changed")
			}
			luogo := htmlquery.InnerText(luogoNode)
			luogo = strings.TrimSpace(luogo)

			exam := Exam{
				SubjectCode:   subjCode,
				SubjectName:   title,
				Teacher:       subjTeacher,
				Date:          date,
				Type:          tipoProva,
				Location:      luogo,
				Subscriptions: listaIscrizioni,
			}

			exams = append(exams, exam)
		}
	}

	return exams, nil
}
