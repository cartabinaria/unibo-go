// SPDX-FileCopyrightText: 2025 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: MIT

package rubrica

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/karlseguin/ccache/v3"
)

type Contact struct {
	FirstName string
	LastName  string
	WorkTitle string
	Phone     string
	Email     string
	WebSite   string
}

type Contacts []Contact

var baseUrl = "https://www.unibo.it/uniboweb/unibosearch/rubrica.aspx?tab=PersonePanel&mode=people&query="

var searchCache = ccache.New(ccache.Configure[Contacts]().MaxSize(1000))

func Search(firstName, lastName string) (Contacts, error) {
	url := baseUrl
	if firstName != "" {
		url += "+nome:" + firstName
	}
	if lastName != "" {
		url += "+cognome:" + lastName
	}

	// try cache
	item := searchCache.Get(url)
	if item != nil {
		return item.Value(), nil
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

	var contacts Contacts
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

	// cache
	searchCache.Set(url, contacts, 0)

	return contacts, nil
}
