package opendata

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/cartabinaria/unibo-go/ckan"
	"github.com/cartabinaria/unibo-go/degree"
)

const (
	// PackageDegreeProgrammesId is the id of the package containing the degrees.
	PackageDegreeProgrammesId = "degree-programmes"

	// ResourceDegreeProgrammesAlias is the alias of the resource containing the
	// degrees.
	ResourceDegreeProgrammesAlias = "corsi_latest_it"
)

// GetDegrees fetches and returns the degrees available in the open data for the
// current year.
//
// Internally it uses PackageDegreeProgrammesId and ResourceDegreeProgrammesAlias
// constants to get the wanted package and resource and then parses the csv file
// to get the degrees.
//
// If you'd like to get the package and resource yourself (maybe you want to
// download the csv file), you can use FetchPackage and Resources.GetByAlias.
func GetDegrees() ([]degree.Degree, error) {
	// Get package
	pack, err := ckanClient.GetPackage(PackageDegreeProgrammesId)
	if err != nil {
		return nil, err
	}

	// If no resources, return nil
	if len(pack.Resources) == 0 {
		return nil, errors.New("no resources found while downloading degrees open data")
	}

	// Get wanted resource
	resource, found := ckan.GetByAlias(pack.Resources, ResourceDegreeProgrammesAlias)
	if !found {
		return nil, errors.New("unable to find resource '" + ResourceDegreeProgrammesAlias + "'")
	}

	// Get the resource
	res, err := http.Get(resource.URL)
	if err != nil {
		return nil, fmt.Errorf("unable to get resource: %w", err)
	} else if res.Header.Get("Content-Type") != "text/csv" {
		return nil, errors.New("resource is not a csv file")
	}

	// Parse the body
	reader := csv.NewReader(res.Body)

	// Skip first line
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("unable to skip first line: %w", err)
	}

	degrees := make([]degree.Degree, 0, 100)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("unable to read row: %w", err)
		}

		if len(row) != 15 {
			return nil, fmt.Errorf("unexpected number of fields: %d (the website has changed?)", len(row))
		}

		degreeCode := row[2]

		years, err := strconv.ParseInt(row[9], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("unable to parse field 'years': %w", err)
		}

		international, err := strconv.ParseBool(row[10])
		if err != nil {
			return nil, fmt.Errorf("unable to parse field 'international': %w", err)
		}

		degrees = append(degrees, degree.Degree{
			AcademicYear:          row[0],
			OpenForRegistration:   row[1],
			Code:                  degreeCode,
			Description:           row[3],
			Url:                   row[4],
			Campus:                row[5],
			TeachingLocation:      row[6],
			Fields:                row[7],
			Type:                  row[8],
			DurationInYears:       int(years),
			International:         international,
			InternationalTitle:    row[11],
			InternationalLanguage: row[12],
			Languages:             row[13],
			AccessRequirements:    row[14],
		})
	}

	return degrees, nil
}
