package timetable

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchTimetable(t *testing.T) {

	const mockTimetable = `[{"cod_modulo":"28004_1","periodo_calendario":"18 settembre 2023 - 20 dicembre 2023","cod_sdoppiamento":"28004_1--A-K","title":"FONDAMENTI DI INFORMATICA T-1 / (A-K) / (1) Modulo 1","periodo":"19 settembre 2023 - 20 dicembre 2023","extCode":"2023-000-356354--I","aule":[{"des_indirizzo":"Viale del Risorgimento, 2 - Bologna","des_piano":"Piano Secondo","des_edificio":"AULA 6.2","raw":{"abilitato":true,"capienzaEffettiva":250,"tipoAulaId":"5dc3ed7874895700123a9c61","unitaOrganizzativaId":"5dc3ed4c74895700123a8b46","attivo":true,"dataCreazione":"2019-11-08T12:51:03.483Z","utenteModificaId":"5dc2de876904b9bd4c70a30b","dataDisabilitazione":"2023-03-02T23:37:56.898Z","dataModifica":"2023-03-06T23:37:50.900Z","bloccato":true,"metriQuadri":218.23000000000002,"piano":{"codice":"331_WP02","descrizione":"Piano Secondo"},"numeroPostazioni":250,"utenteCreazioneId":"5dc2de876904b9bd4c70a30b","extCode":"331_WP02_234","utenteDisabilitazioneId":"5dc2de876904b9bd4c70a30b","edificio":{"comune":"Bologna","via":"Viale del Risorgimento, 2","provincia":"Bologna","codice":"331","utenteCreazioneId":"5dc2de876904b9bd4c70a30b","attivo":true,"extCode":"331","cap":"40136","descrizione":"Facoltà di Ingegneria dell'Università di Bologna","plesso":"RISORGIMENTO","dataModifica":"2023-07-20T22:35:22.041Z","dataCreazione":"2019-11-07T10:09:32.720Z","metriQuadri":28760.79,"utenteModificaId":"5dc2de876904b9bd4c70a30b","clienteId":"5ad08435b6ca5357dbac609e","_id":"5dc3ed5c74895700123a91aa","geo":{"lat":44.4903628,"lng":11.3289228},"bloccato":true},"capienza":250,"mappaUrl":null,"clienteId":"5ad08435b6ca5357dbac609e","fotoUrl":null,"serviziAulaId":[],"_id":"5dc564b7b2285f0011f80c4b","codice":"331_WP02_234","divisoreCapienza":1,"planimetriaUrl":null,"tipoAula":{"codice":"AULA","descrizione":"AULA"},"descrizione":"AULA 6.2","tolleranza":0,"edificioId":"5dc3ed5c74895700123a91aa","pianoId":"5dc3ed6574895700123a95e1"},"des_ubicazione":"Facoltà di Ingegneria dell'Università di Bologna","des_risorsa":"AULA 6.2"}],"teledidattica":false,"teams":null,"note":"","start":"2023-09-19T09:00:00","docente":"Paola Mello","time":"09:00 - 12:00","end":"2023-09-19T12:00:00","cfu":12,"val_crediti":12}]`

	// Mock the HTTP client
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/laurea/IngegneriaInformatica/orario-lezioni/@@orario_reale_json" {
			t.Error("wrong path", r.URL.Path)
		}
		if r.URL.RawQuery != "anno=1" {
			t.Error("wrong query", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(mockTimetable))
	}))
	defer server.Close()

	baseUrl = server.URL
	http.DefaultClient = server.Client()

	// Fetch the timetable
	timetable, err := FetchTimetable("laurea", "IngegneriaInformatica", "", 1, nil)
	if err != nil {
		t.Error(err)
	}

	if len(timetable) != 1 {
		t.Error("wrong number of events", len(timetable))
	}

	if timetable[0].CodModulo != "28004_1" {
		t.Error("wrong CodModulo", timetable[0].CodModulo)
	}

	if timetable[0].CalendarInterval != "18 settembre 2023 - 20 dicembre 2023" {
		t.Error("wrong CalendarInterval")
	}

	if timetable[0].CodSdoppiamento != "28004_1--A-K" {
		t.Error("wrong CodSdoppiamento")
	}

	if timetable[0].Title != "FONDAMENTI DI INFORMATICA T-1 / (A-K) / (1) Modulo 1" {
		t.Error("wrong Title")
	}

	if timetable[0].ExtCode != "2023-000-356354--I" {
		t.Error("wrong ExtCode")
	}

	if timetable[0].Interval != "19 settembre 2023 - 20 dicembre 2023" {
		t.Error("wrong Interval")
	}

	if timetable[0].Teacher != "Paola Mello" {
		t.Error("wrong Teacher")
	}

	if timetable[0].Cfu != 12 {
		t.Error("wrong Cfu")
	}

	if timetable[0].RemoteLearning != false {
		t.Error("wrong RemoteLearning")
	}

	if timetable[0].Teams != "" {
		t.Error("wrong Teams")
	}

	if timetable[0].Start.String() != "2023-09-19 09:00:00 +0200 CEST" {
		t.Error("wrong Start", timetable[0].Start.String())
	}

	if timetable[0].End.String() != "2023-09-19 12:00:00 +0200 CEST" {
		t.Error("wrong End", timetable[0].End.String())
	}

	if len(timetable[0].Classrooms) != 1 {
		t.Error("wrong number of classrooms")
	} else if timetable[0].Classrooms[0].ResourceDesc != "AULA 6.2" {
		t.Error("wrong Description")
	}
}
