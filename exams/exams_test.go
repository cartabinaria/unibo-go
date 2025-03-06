package exams

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetExams(t *testing.T) {
	// Mock server to simulate the Unibo website
	handler := http.NewServeMux()
	handler.HandleFunc("/unibo/test/appelli", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("b_start:int") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
<div class=dropdown-component role=tablist>
    <h3 class="border-secondary background-secondary-dark" aria-controls=panel0 aria-expanded=true aria-selected=true
        id=tab0 role=tab><a href=# class=openclose-appelli><span class=code>72677</span> ANALISI DELLE RETI SOCIALI
            APPLICATA AD INTERNET <span class=docente>GIALLORENZO SAVERIO</span> <i aria-hidden=true
                class="fa fa-caret-up"></i></a></h3>
    <div class=items-container role=tabpanel aria-hidden=false aria-labelledby=tab0 id=panel0>
        <table class=single-item>
            <tr>
                <th class=text-secondary>Data e ora:
                <td class=text-secondary>06 dicembre 2024 ore 09:00
            <tr>
                <th>Lista iscrizioni:
                <td>aperta dal <span>18 ottobre 2024</span> al <span>05 dicembre 2024</span>
            <tr>
                <th>Tipo prova:
                <td>Scritto
            <tr>
                <th>Luogo:
                <td>ONLINE
        </table>
        <table class=single-item>
            <tr>
                <th class=text-secondary>Data e ora:
                <td class=text-secondary>10 gennaio 2025 ore 09:00
            <tr>
                <th>Lista iscrizioni:
                <td>aperta dal <span>26 dicembre 2024</span> al <span>09 gennaio 2025</span>
            <tr>
                <th>Tipo prova:
                <td>Scritto
            <tr>
                <th>Luogo:
                <td>ONLINE
        </table>
        <table class=single-item>
            <tr>
                <th class=text-secondary>Data e ora:
                <td class=text-secondary>21 febbraio 2025 ore 09:00
            <tr>
                <th>Lista iscrizioni:
                <td>aperta dal <span>13 gennaio 2025</span> al <span>20 febbraio 2025</span>
            <tr>
                <th>Tipo prova:
                <td>Scritto
            <tr>
                <th>Luogo:
                <td>ONLINE
        </table>
        <table class=single-item>
            <tr>
                <th class=text-secondary>Data e ora:
                <td class=text-secondary>04 aprile 2025 ore 09:00
            <tr>
                <th>Lista iscrizioni:
                <td>aperta dal <span>20 marzo 2025</span> al <span>03 aprile 2025</span>
            <tr>
                <th>Tipo prova:
                <td>Scritto
            <tr>
                <th>Luogo:
                <td>ONLINE
        </table>
    </div>
</div>
		`))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	// Override the baseUrl to point to the mock server
	baseUrl = server.URL + "/%s/%s/appelli"

	exams, err := GetExams("unibo", "test")
	assert.NoError(t, err)
	assert.NotNil(t, exams)
	assert.Equal(t, 4, len(exams))

	assert.Equal(t, "72677", exams[0].SubjectCode)
	assert.Equal(t, "ANALISI DELLE RETI SOCIALI APPLICATA AD INTERNET", exams[0].SubjectName)
	assert.Equal(t, time.Date(2024, time.December, 06, 9, 0, 0, 0, time.UTC), exams[0].Date)
	assert.Equal(t, "aperta dal 18 ottobre 2024 al 05 dicembre 2024", exams[0].Subscriptions)
	assert.Equal(t, "Scritto", exams[0].Type)
	assert.Equal(t, "ONLINE", exams[0].Location)

	assert.Equal(t, "72677", exams[1].SubjectCode)
	assert.Equal(t, "ANALISI DELLE RETI SOCIALI APPLICATA AD INTERNET", exams[1].SubjectName)
	assert.Equal(t, time.Date(2025, time.January, 10, 9, 0, 0, 0, time.UTC), exams[1].Date)
	assert.Equal(t, "aperta dal 26 dicembre 2024 al 09 gennaio 2025", exams[1].Subscriptions)
	assert.Equal(t, "Scritto", exams[1].Type)
	assert.Equal(t, "ONLINE", exams[1].Location)

}
