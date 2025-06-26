// SPDX-FileCopyrightText: 2025 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: MIT

package rubrica

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
)

type Contact struct {
	FirstName string
	LastName  string
	WorkTitle string
	Phone     string
	Email     string
	WebSite   string
}

// These are declared as variables to allow for easier testing and mocking
var (
	baseUrl    = "https://www.unibo.it/uniboweb/unibosearch/rubrica.aspx?tab=PersonePanel&mode=people&query="
	httpClient = &http.Client{}
)

func Search(firstName, lastName string) ([]Contact, error) {
	url := baseUrl
	if firstName != "" {
		url += "+nome:" + firstName
	}
	if lastName != "" {
		url += "+cognome:" + lastName
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to get res: %w", err)
	}

	// parse res
	node, err := htmlquery.Parse(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to parse res: %w", err)
	}

	tables := htmlquery.Find(node, "//table[@class='contact vcard']")
	if len(tables) == 0 {
		return nil, nil
	}

	var contacts []Contact
	for _, table := range tables {
		fullName := htmlquery.FindOne(table, "//td[@class='fn name']")
		// split on ,
		names := strings.Split(htmlquery.InnerText(fullName), ",")
		if len(names) != 2 {
			return nil, fmt.Errorf("unable to split name: %s", htmlquery.InnerText(fullName))
		}

		firstName := strings.TrimSpace(names[1])
		lastName := strings.TrimSpace(names[0])

		emailNode := htmlquery.InnerText(htmlquery.FindOne(table, "//a[@class='email']"))
		email := strings.TrimPrefix(emailNode, "mailto:")

		contacts = append(contacts, Contact{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
		})

	}

	return contacts, nil
}
