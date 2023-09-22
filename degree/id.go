package degree

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type ID struct {
	Type string
	Id   string
}

var reg = regexp.MustCompile(`<a title="Sito del corso" href="https://corsi\.unibo\.it/(.+?)"`)

// ScrapeId returns the ID of the course from the given course website url.
func (d Degree) ScrapeId() (ID, error) {

	resp, err := http.Get(d.Url)
	if err != nil {
		return ID{}, err
	}

	buf := new(bytes.Buffer)

	// Read all body
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return ID{}, err
	}

	// Close body
	err = resp.Body.Close()
	if err != nil {
		return ID{}, err
	}

	// Convert body to string
	found := reg.FindStringSubmatch(buf.String())
	if found == nil {
		return ID{}, fmt.Errorf("unable to find course website")
	} else if len(found) != 2 {
		return ID{}, fmt.Errorf("unexpected number of matches: %d (the website has changed?)", len(found))
	}

	// full url -> laurea/IngegneriaInformatica
	id := found[1]

	// laurea/IngegneriaInformatica -> IngegneriaInformatica
	split := strings.Split(id, "/")
	return ID{Type: split[0], Id: split[1]}, nil
}
