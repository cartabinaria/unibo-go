package ckan

import (
	"slices"
	"strings"
)

// GetByAlias returns the resource with the given alias, if it exists in the Resources slice.
//
// If the resource is found, it returns a pointer to the resource and true.
// Otherwise, it returns nil and false.
func GetByAlias(r []Resource, alias string) (*Resource, bool) {
	for _, resource := range r {
		// Some resources have multiple aliases
		rAliases := strings.Split(resource.Alias, ", ")

		// Check if the alias is one of the aliases of the resource
		if slices.Contains(rAliases, alias) {
			return &resource, true
		}
	}
	return nil, false
}
