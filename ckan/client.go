package ckan

import "fmt"

// Client represents a CKAN client that can interact with a CKAN instance.
//
// It provides methods to fetch data such as packages, groups, organizations,
// tags, users, licenses, and vocabularies from the CKAN API.
type Client struct {
	baseURL string // The base URL of the CKAN instance
}

// NewClient creates a new CKAN client with the given base URL.
func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

func (c *Client) GetPackageList() (*[]string, error) {
	return Request[[]string](fmt.Sprintf("%s/api/3/action/package_list", c.baseURL))
}
func (c *Client) GetPackage(id string) (*Package, error) {
	return Request[Package](fmt.Sprintf("%s/api/3/action/package_show?id=%s", c.baseURL, id))
}
func (c *Client) GetCurrentPackageListWithResources(limit, offset int) (*[]Package, error) {
	return Request[[]Package](fmt.Sprintf("%s/api/3/action/current_package_list_with_resources?limit=%d&offset=%d", c.baseURL, limit, offset))
}
func (c *Client) GetGroupList() (*[]string, error) {
	return Request[[]string](fmt.Sprintf("%s/api/3/action/group_list", c.baseURL))
}
func (c *Client) GetGroup(id string) (*Group, error) {
	return Request[Group](fmt.Sprintf("%s/api/3/action/group_show?id=%s", c.baseURL, id))
}
func (c *Client) GetOrganizationList() (*[]string, error) {
	return Request[[]string](fmt.Sprintf("%s/api/3/action/organization_list", c.baseURL))
}
func (c *Client) GetOrganization(id string) (*Organization, error) {
	return Request[Organization](fmt.Sprintf("%s/api/3/action/organization_show?id=%s", c.baseURL, id))
}
func (c *Client) GetTagList() (*[]string, error) {
	return Request[[]string](fmt.Sprintf("%s/api/3/action/tag_list", c.baseURL))
}
func (c *Client) GetTagShow(id string) (*Tag, error) {
	return Request[Tag](fmt.Sprintf("%s/api/3/action/tag_show?id=%s", c.baseURL, id))
}
func (c *Client) GetUserList() (*[]string, error) {
	return Request[[]string](fmt.Sprintf("%s/api/3/action/user_list", c.baseURL))
}
func (c *Client) GetUser(id string) (*User, error) {
	return Request[User](fmt.Sprintf("%s/api/3/action/user_show?id=%s", c.baseURL, id))
}
func (c *Client) GetLicenseList() (*[]License, error) {
	return Request[[]License](fmt.Sprintf("%s/api/3/action/license_list", c.baseURL))
}
func (c *Client) GetVocabularyList() (*[]Vocabulary, error) {
	return Request[[]Vocabulary](fmt.Sprintf("%s/api/3/action/vocabulary_list", c.baseURL))
}
func (c *Client) GetVocabulary(id string) (*Vocabulary, error) {
	return Request[Vocabulary](fmt.Sprintf("%s/api/3/action/vocabulary_show?id=%s", c.baseURL, id))
}
func (c *Client) GetPackageSearch(query string, rows, start int) (*PackageSearch, error) {
	return Request[PackageSearch](fmt.Sprintf("%s/api/3/action/package_search?q=%s&rows=%d&start=%d", c.baseURL, query, rows, start))
}
func (c *Client) SearchResource(query string, limit, offset int) (*ResourceSearch, error) {
	return Request[ResourceSearch](fmt.Sprintf("%s/api/3/action/resource_search?query=%s&limit=%d&offset=%d", c.baseURL, query, limit, offset))
}
func (c *Client) SearchTag(query string, limit, offset int) (*TagSearch, error) {
	return Request[TagSearch](fmt.Sprintf("%s/api/3/action/tag_search?query=%s&limit=%d&offset=%d", c.baseURL, query, limit, offset))
}
