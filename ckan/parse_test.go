package ckan

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func serveJson(t *testing.T, query, json string) string {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		completePath := r.URL.Path + "?" + r.URL.RawQuery
		require.Equal(t, query, completePath, "Unexpected request path")

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(json))
		require.NoError(t, err, "Failed to write response")
	}))

	t.Cleanup(func() {
		srv.Close()
	})

	return srv.URL
}

func TestGetGroupShow(t *testing.T) {
	json := `{"help": "https://demo.ckan.org/api/3/action/help_show?name=group_show", "success": true, "result": {"approval_status": "approved", "created": "2024-02-27T14:02:47.679370", "description": "These are books that David likes.", "display_name": "Dave's books", "id": "ded220b4-c665-4312-95d0-8c3ec969441f", "image_display_url": "", "image_url": "", "is_organization": false, "name": "david", "num_followers": 0, "package_count": 2, "state": "active", "title": "Dave's books", "type": "group", "extras": [], "tags": [], "groups": []}}`
	url := serveJson(t, "/api/3/action/group_show?id=david", json)

	client := NewClient(url)
	groups, err := client.GetGroup("david")

	require.NoError(t, err, "GetGroup should not return an error")
	require.NotNil(t, groups, "GetGroup should return a non-nil group")

	assert.Equal(t, "ded220b4-c665-4312-95d0-8c3ec969441f", groups.ID, "Expected group ID to match")
	assert.Equal(t, "david", groups.Name, "Expected group name to match")
	assert.Equal(t, "These are books that David likes.", groups.Description, "Expected group description to match")
	assert.Equal(t, "Dave's books", groups.Title, "Expected group title to match")
}

//go:embed testdata/package.json
var getPackageTestData string

func TestGetPackage(t *testing.T) {

	url := serveJson(t, "/api/3/action/package_show?id=1", getPackageTestData)
	client := NewClient(url)

	pkg, err := client.GetPackage("1")
	require.NoError(t, err, "GetResource should not return an error")
	require.NotNil(t, pkg, "GetResource should return a non-nil package")

	assert.Equal(t, "8cbba6bd-1686-4828-b223-6f4dab477434", pkg.ID, "Expected package ID to match")
	assert.Equal(t, "Corsi di studio", pkg.Title, "Expected package title to match")
	assert.Equal(t, "degree-programmes", pkg.Name, "Expected package name to match")
	assert.Equal(t, "active", pkg.State, "Expected package state to match")
	assert.Equal(t, "dataset", pkg.Type, "Expected package type to match")
	assert.Equal(t, "CC-BY-3.0-IT", pkg.LicenseID, "Expected license ID to match")
	assert.Equal(t, "opendata@unibo.it", pkg.MaintainerEmail, "Expected maintainer email to match")
	assert.Equal(t, false, pkg.Private, "Expected package to be public")
	assert.Equal(t, 0, pkg.NumTags, "Expected number of tags to match")
	assert.Equal(t, "2.3.1", pkg.Version, "Expected version to match")
	assert.Equal(t, "documentation-degree-programmes", pkg.Documentation, "Expected documentation to match")
	assert.Contains(t, pkg.Notes, "Catalogo dei corsi di Laurea", "Expected notes to contain Italian description")
	assert.Contains(t, pkg.NotesTranslated["en"], "Catalogue of the First and Second cycle degree programmes", "Expected English notes to match")
	assert.Contains(t, pkg.NotesTranslated["it"], "Catalogo dei corsi di Laurea", "Expected Italian notes to match")
	assert.Contains(t, pkg.TitleTranslated["en"], "Degree programmes", "Expected English title to match")
	assert.Contains(t, pkg.TitleTranslated["it"], "Corsi di studio", "Expected Italian title to match")
	assert.NotEmpty(t, pkg.Resources, "Expected resources to be present")
	assert.GreaterOrEqual(t, len(pkg.Resources), 1, "Expected at least one resource")
	assert.Equal(t, "d47703c3-aba7-4e69-9a31-313c4c3ffcb0", pkg.Resources[0].ID, "Expected first resource ID to match")
	assert.Equal(t, "Corsi di studio 2024/2025", pkg.Resources[0].Name, "Expected first resource name to match")
	assert.Equal(t, "CSV", pkg.Resources[0].Format, "Expected first resource format to match")
	assert.Equal(t, "active", pkg.Resources[0].State, "Expected first resource state to match")
	assert.Equal(t, "https://dati.unibo.it/dataset/8cbba6bd-1686-4828-b223-6f4dab477434/resource/d47703c3-aba7-4e69-9a31-313c4c3ffcb0/download/corsi_2024_it.csv", pkg.Resources[0].URL, "Expected first resource URL to match")
	assert.Equal(t, "it", pkg.Resources[0].Language, "Expected first resource language to match")
	assert.Equal(t, "DAILY", pkg.Resources[0].Frequency, "Expected first resource frequency to match")
}

//go:embed testdata/organization.json
var getOrganizationTestData string

func TestGetOrganization(t *testing.T) {
	srv := serveJson(t, "/api/3/action/organization_show?id=aform", getOrganizationTestData)
	client := NewClient(srv)
	org, err := client.GetOrganization("aform")

	require.NoError(t, err, "GetOrganization should not return an error")
	require.NotNil(t, org, "GetOrganization should return a non-nil organization")

	assert.Equal(t, "31f7674e-2f06-4b54-a833-d630bf8463ff", org.ID, "Expected organization ID to match")
	assert.Equal(t, "aform", org.Name, "Expected organization name to match")
	assert.Equal(t, "AFORM - Education Division", org.DisplayName, "Expected display name to match")
	assert.Equal(t, " AFORM - Area della Didattica", org.Title, "Expected title to match")
	assert.Equal(t, "organization", org.Type, "Expected type to match")
	assert.Equal(t, "active", org.State, "Expected state to match")
	assert.Equal(t, 26, org.PackageCount, "Expected package count to match")
	assert.Equal(t, "approved", org.ApprovalStatus, "Expected approval_status to match")
	assert.Equal(t, "", org.Description, "Expected description to be empty")
	assert.Equal(t, "", org.Notes, "Expected notes to be empty")
	assert.Equal(t, "", org.ImageDisplayURL, "Expected image_display_url to be empty")
	assert.Equal(t, "", org.ImageURL, "Expected image_url to be empty")
	assert.Equal(t, 0, org.NumFollowers, "Expected num_followers to be 0")
	assert.Equal(t, "b3f3b914-2c0c-460d-aa17-3d8342897236", org.RevisionID, "Expected revision_id to match")
	assert.Equal(t, "2016-12-02T16:35:27.570738", org.Created, "Expected created date to match")
	assert.Equal(t, map[string]string{"en": "AFORM - Education Division", "it": " AFORM - Area della Didattica"}, org.TitleTranslated, "Expected title_translated to match")
	assert.Equal(t, TranslatedValue{"en": "", "it": ""}, org.NotesTranslated, "Expected notes_translated to match")
	require.Len(t, org.Users, 5, "Expected 5 users in organization")
	assert.Equal(t, "e2101ecc-e65f-4f11-bb3d-f2474673fecc", org.Users[0].ID, "Expected first user ID to match")
	assert.Equal(t, "admin", org.Users[0].Name, "Expected first user name to match")
	assert.Equal(t, "admin", org.Users[0].DisplayName, "Expected first user display name to match")
	assert.Equal(t, "admin", org.Users[0].Capacity, "Expected first user capacity to match")
	assert.Equal(t, true, org.Users[0].Sysadmin, "Expected first user sysadmin to be true")
	assert.Equal(t, "active", org.Users[0].State, "Expected first user state to match")
	assert.Equal(t, 67375, org.Users[0].NumberOfEdits, "Expected first user number_of_edits to match")
	assert.Equal(t, 5, org.Users[0].NumberCreatedPackages, "Expected first user number_created_packages to match")
	assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", org.Users[0].EmailHash, "Expected first user email_hash to match")
	assert.Equal(t, "2016-11-23T13:09:37.301082", org.Users[0].Created, "Expected first user created date to match")
	assert.Equal(t, "01d2fd29-f1a5-4b18-9971-873255f2cef5", org.Users[1].ID, "Expected second user ID to match")
	assert.Equal(t, "carlo-marchesi-7979", org.Users[1].Name, "Expected second user name to match")
	assert.Equal(t, "pending", org.Users[1].State, "Expected second user state to match")
	assert.Equal(t, "member", org.Users[1].Capacity, "Expected second user capacity to match")
	assert.Equal(t, false, org.Users[1].Sysadmin, "Expected second user sysadmin to be false")
	assert.Equal(t, "4b2bc546-c679-495a-a980-1b4d806f2d3d", org.Users[2].ID, "Expected third user ID to match")
	assert.Equal(t, "carlo.marchesi@unibo.it", org.Users[2].Name, "Expected third user name to match")
	assert.Equal(t, "active", org.Users[2].State, "Expected third user state to match")
	assert.Equal(t, "member", org.Users[2].Capacity, "Expected third user capacity to match")
	assert.Equal(t, false, org.Users[2].Sysadmin, "Expected third user sysadmin to be false")
	assert.Equal(t, "9c06c8e9-9274-4703-8342-5d1831542220", org.Users[3].ID, "Expected fourth user ID to match")
	assert.Equal(t, "serena-alessandrini-7365", org.Users[3].Name, "Expected fourth user name to match")
	assert.Equal(t, "pending", org.Users[3].State, "Expected fourth user state to match")
	assert.Equal(t, "member", org.Users[3].Capacity, "Expected fourth user capacity to match")
	assert.Equal(t, false, org.Users[3].Sysadmin, "Expected fourth user sysadmin to be false")
	assert.Equal(t, "64a43a5e-03e9-4930-a492-a58bebe80cad", org.Users[4].ID, "Expected fifth user ID to match")
	assert.Equal(t, "serena.alessandrini@unibo.it", org.Users[4].Name, "Expected fifth user name to match")
	assert.Equal(t, "active", org.Users[4].State, "Expected fifth user state to match")
	assert.Equal(t, "member", org.Users[4].Capacity, "Expected fifth user capacity to match")
	assert.Equal(t, false, org.Users[4].Sysadmin, "Expected fifth user sysadmin to be false")
}
