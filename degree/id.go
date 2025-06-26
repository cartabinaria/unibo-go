// SPDX-FileCopyrightText: 2023 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: MIT

package degree

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// ID represents the ID of a course. It is made of a type and an id.
type ID struct {
	Type string // Type is the type of the course, e.g. "laurea".
	Id   string // Id is the id of the course, e.g. "IngegneriaInformatica".
}

var reg = regexp.MustCompile(`<a title="Sito del corso" href="https://corsi\.unibo\.it/(.+?)"`)

// ScrapeId returns the ID of the course from the given course website url.
func (d *Degree) ScrapeId() (ID, error) {

	resp, err := http.Get(d.Url)
	if err != nil {
		return ID{}, fmt.Errorf("could not get course website: %w", err)
	}

	// Read all body in memory
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return ID{}, fmt.Errorf("could not read response body: %w", err)
	}

	// Close body
	err = resp.Body.Close()
	if err != nil {
		return ID{}, fmt.Errorf("could not close response body: %w", err)
	}

	// Convert body to string
	found := reg.FindStringSubmatch(string(buf))
	if found == nil {
		return ID{}, errors.New("unable to find course website")
	} else if len(found) != 2 {
		return ID{}, fmt.Errorf("unexpected number of matches: %d (the website has changed?)", len(found))
	}

	// full url -> laurea/IngegneriaInformatica
	id := found[1]

	// laurea/IngegneriaInformatica -> IngegneriaInformatica
	split := strings.Split(id, "/")
	return ID{Type: split[0], Id: split[1]}, nil
}
