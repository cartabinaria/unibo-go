// Package opendata provides functions to fetch data from the UniBo Open Data
// portal.
//
// It is useful in particular to get degrees (see GetDegrees).
package opendata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Resource represents a resource in the UniBO Open Data portal.
type Resource struct {
	Frequency string `json:"frequency"`
	Url       string `json:"url"`
	Id        string `json:"id"`
	PackageId string `json:"package_id"`
	LastMod   string `json:"last_modified"`
	Alias     string `json:"alias"`
}
type Resources []Resource

// Package represents a package in the UniBO Open Data portal.
type Package struct {
	Success bool `json:"success"`
	Result  struct {
		Resources Resources
	}
}

// GetByAlias returns the resource with the given alias, if it exists in the Resources slice.
//
// If the resource is found, it returns a pointer to the resource and true.
// Otherwise, it returns nil and false.
func (r Resources) GetByAlias(alias string) (*Resource, bool) {
	for _, resource := range r {
		// Some resources have multiple aliases
		rAliases := strings.Split(resource.Alias, ", ")

		// Check if the alias is one of the aliases of the resource
		for _, rAlias := range rAliases {
			if rAlias == alias {
				return &resource, true
			}
		}
	}
	return nil, false
}

const openDataUrl = "https://dati.unibo.it"

// getPackageUrl returns the url to fetch the Package with the given id.
func getPackageUrl(id string) string {
	return fmt.Sprintf("%s/api/3/action/package_show?id=%s", openDataUrl, id)
}

// FetchPackage retrieves the package with the given id from the Unibo Open Data portal.
//
// If the package is found, it returns a pointer to the package.
// Otherwise, it returns nil and an error.
func FetchPackage(id string) (*Package, error) {
	url := getPackageUrl(id)

	html, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if html.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", html.StatusCode, html.Status)
	}

	body := html.Body
	pack := Package{}

	err = json.NewDecoder(body).Decode(&pack)
	if err != nil {
		return nil, err
	}

	err = body.Close()
	if err != nil {
		return nil, err
	}

	return &pack, nil
}
