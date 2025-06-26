// SPDX-FileCopyrightText: 2025 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: MIT

// The CKAN package provides functionality to interact with the CKAN API,
// a popular open data portal software.
//
// To use this package, first create a new client with the base URL of your CKAN instance.
//
//	client := ckan.NewClient("https://your-ckan-instance.com")
//
// You can then use the client to fetch various resources from the CKAN API.
package ckan

// KVPair is a generic key-value pair.
type KVPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TranslatedValue is a map that holds translations for a value in different languages.
//
// Keys are language codes (e.g., "en", "it"), and values are the translated strings.
type TranslatedValue map[string]string

// Activity represents an activity in CKAN, which can be related to a Package, Group, or User.
type Activity struct {
	ID           string         `json:"id"`            // The unique identifier for the activity
	ObjectID     string         `json:"object_id"`     // The ID of the object related to the activity, which can be a Package, Group, or User
	ObjectType   string         `json:"object_type"`   // The type of the object related to the activity, e.g., "package", "group", "user"
	ActivityType string         `json:"activity_type"` // The type of activity, e.g., "create", "update", "delete"
	Timestamp    string         `json:"timestamp"`     // The date and time when the activity occurred
	UserID       string         `json:"user_id"`       // The ID of the user who performed the activity
	Data         map[string]any `json:"data"`          // Additional data related to the activity, which can vary based on the activity type
}

// Group represents a CKAN group, which is a collection of Packages.
type Group struct {
	ID              string   `json:"id"`                // The ID of the group
	Name            string   `json:"name"`              // The name of the group
	Title           string   `json:"title"`             // The title of the group
	Description     string   `json:"description"`       // A description of the group
	ImageDisplayURL string   `json:"image_display_url"` // The URL of the group's display image. Can be empty if no image is set.
	Created         string   `json:"created"`           // The date and time when the group was created
	ApprovalStatus  string   `json:"approval_status"`   // The approval status of the group, e.g., "approved", "pending"
	State           string   `json:"state"`             // The state of the group, e.g., "active", "deleted"
	Type            string   `json:"type"`              //
	Users           []User   `json:"users"`             // The list of users associated with the group
	Packages        []string `json:"packages"`          // The list of package IDs in the group
	Extras          []KVPair `json:"extras"`            // Additional key-value pairs associated with the group
}

// Resource represents a CKAN resource, which is a file or data source associated with a Package.
//
// A resource can be a file, a URL, or any other data source that provides information related to the Package.
type Resource struct {
	ID               string  `json:"id"`                 // The ID of the resource
	PackageID        string  `json:"package_id"`         // The ID of the package this resource belongs to
	URL              string  `json:"url"`                // The URL of the resource, if applicable
	Name             string  `json:"name"`               // The name of the resource. Can be empty if not set.
	Description      string  `json:"description"`        // A description of the resource
	Format           string  `json:"format"`             // The format of the resource, e.g., "CSV", "JSON"
	Mimetype         string  `json:"mimetype"`           // The MIME type of the resource, can be nil
	MimetypeInner    string  `json:"mimetype_inner"`     //
	Created          string  `json:"created"`            // The date and time when the resource was created
	LastModified     string  `json:"last_modified"`      // The date and time when the resource was last modified. Can be empty if not set.
	MetadataModified string  `json:"metadata_modified"`  //
	CacheURL         string  `json:"cache_url"`          // The URL to the cached version of the resource, if applicable.
	CacheLastUpdated string  `json:"cache_last_updated"` // The date and time when the cache was last updated. Can be empty if not set.
	Size             int64   `json:"size"`               // The size of the resource in bytes.
	SizeExtra        string  `json:"size_extra"`         //
	State            string  `json:"state"`              //
	Hash             string  `json:"hash"`               //
	Position         int     `json:"position"`           //
	ResourceType     *string `json:"resource_type"`      //
	RevisionID       *string `json:"revision_id"`        //
	DatastoreActive  bool    `json:"datastore_active"`   //
	URLType          *string `json:"url_type"`           //
	Alias            string  `json:"alias"`              // Some resources have multiple aliases
	Language         string  `json:"language"`           // The language of the resource, e.g., "en", "it"
	Frequency        string  `json:"frequency"`          // The frequency of updates for the resource, e.g., "DAILY", "WEEKLY"
}

// Package represents a CKAN package, which is a collection of resources and metadata.
//
// A package can be thought of as a dataset or a collection of related data.
type Package struct {
	ID               string          `json:"id"`               // The unique identifier for the package
	Name             string          `json:"name"`             // The name of the package, typically a URL-friendly version of the title
	Title            string          `json:"title"`            // The title of the package
	TitleTranslated  TranslatedValue `json:"title_translated"` // Translated titles for the package in different languages
	Notes            string          `json:"notes"`            // The notes or description of the package
	NotesTranslated  TranslatedValue `json:"notes_translated"` // Translated notes for the package in different languages
	Author           string          `json:"author"`           // The name of the author of the package
	AuthorEmail      string          `json:"author_email"`     // The email of the author of the package
	Maintainer       string          `json:"maintainer"`       // The name of the maintainer of the package
	MaintainerEmail  string          `json:"maintainer_email"` // The email of the maintainer of the package
	LicenseID        string          `json:"license_id"`       // The ID of the license under which the package is released
	LicenseTitle     string          `json:"license_title"`    // The title of the license under which the package is released
	LicenseURL       string          `json:"license_url"`      // The URL of the license under which the package is released
	OwnerOrg         string          `json:"owner_org"`        // The ID of the organization that owns the package
	Organization     *Organization   `json:"organization"`
	Groups           []Group         `json:"groups"`
	Tags             []Tag           `json:"tags"`
	Resources        []Resource      `json:"resources"`
	Extras           []KVPair        `json:"extras"`
	NumResources     int             `json:"num_resources"`
	NumTags          int             `json:"num_tags"`
	State            string          `json:"state"`
	RevisionID       string          `json:"revision_id"`
	MetadataCreated  string          `json:"metadata_created"`
	MetadataModified string          `json:"metadata_modified"`
	Version          string          `json:"version"`
	URL              string          `json:"url"`
	Type             string          `json:"type"`
	IsOpen           bool            `json:"isopen"`
	Private          bool            `json:"private"`       // Indicates if the package is private
	Documentation    string          `json:"documentation"` // URL to the documentation for the package. Can be empty.
}

// Tag represents a CKAN tag, which is a keyword or label used to categorize packages.
//
// Tags can be used to filter and search for packages based on specific topics or themes.
type Tag struct {
	ID                string `json:"id"`                 // The unique identifier for the tag
	Name              string `json:"name"`               // The name of the tag, typically a URL-friendly version of the tag
	DisplayName       string `json:"display_name"`       // The display name of the tag, which can be more user-friendly than the name
	State             string `json:"state"`              // The state of the tag, e.g., "active", "deleted"
	VocabularyID      string `json:"vocabulary_id"`      // The ID of the vocabulary to which the tag belongs, if applicable
	RevisionTimestamp string `json:"revision_timestamp"` // The timestamp of the last revision of the tag
}

// Organization represents a CKAN organization, which is a collection of packages and groups.
//
// An organization can be thought of as a group or entity that manages a set of related datasets.
type Organization struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	DisplayName     string            `json:"display_name"` // The display name of the organization
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	ImageDisplayURL string            `json:"image_display_url"`
	Created         string            `json:"created"`
	ApprovalStatus  string            `json:"approval_status"`
	State           string            `json:"state"`
	Type            string            `json:"type"`
	Extras          []KVPair          `json:"extras"`
	Users           []User            `json:"users"`
	Packages        []string          `json:"packages"`
	PackageCount    int               `json:"package_count"`      // The number of packages in the organization
	Notes           string            `json:"notes"`              // Notes about the organization
	NotesTranslated TranslatedValue   `json:"notes_translated"`   // Translated notes for the organization in different languages
	ImageURL        string            `json:"image_url"`          // The URL of the organization's image, if set
	NumFollowers    int               `json:"num_followers"`      // The number of followers of the organization
	RevisionID      string            `json:"revision_id"`        // The ID of the last revision of the organization
	CreatedByUserID string            `json:"created_by_user_id"` // The ID of the user who created the organization
	TitleTranslated map[string]string `json:"title_translated"`   // Translated titles for the organization in different languages
}

// User represents a CKAN user, which can be an individual or an organization
// that creates and manages packages.
//
// A user can have various roles and permissions within the CKAN instance,
// such as creating packages, managing resources, and participating in groups.
type User struct {
	ID                    string `json:"id"`                      // The unique identifier for the user
	Name                  string `json:"name"`                    // The username of the user
	Fullname              string `json:"fullname"`                // The full name of the user
	DisplayName           string `json:"display_name"`            // The display name of the user
	Email                 string `json:"email"`                   // The email address of the user
	EmailHash             string `json:"email_hash"`              // The MD5 hash of the user's email address
	Created               string `json:"created"`                 // The date and time when the user was created
	NumberCreatedPackages int    `json:"number_created_packages"` // The number of packages created by the user
	About                 string `json:"about"`                   // A short description or bio of the user
	State                 string `json:"state"`                   // The state of the user, e.g., "active", "deleted"
	Sysadmin              bool   `json:"sysadmin"`                // Whether the user is a system administrator
	Capacity              string `json:"capacity"`                // The role of the user in the organization, e.g., "admin", "member"
	NumberOfEdits         int    `json:"number_of_edits"`         // The number of edits made by the user
}

// License represents a CKAN license, which defines the terms under which a package can be used.
//
// A license can be a standard open data license or a custom license defined by the organization or user.
type License struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Vocabulary represents a CKAN vocabulary, which is a collection of tags that
// can be used to categorize packages.
//
// A vocabulary can be thought of as a controlled vocabulary or taxonomy that
// provides a consistent set of tags for use across packages.
type Vocabulary struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Tags   []Tag    `json:"tags"`
	State  string   `json:"state"`
	Extras []KVPair `json:"extras"`
}

// PackageSearch represents the result of a package search query in CKAN.
type PackageSearch struct {
	Count        int            `json:"count"`         // The total number of packages matching the search criteria
	Results      []Package      `json:"results"`       // The list of packages returned by the search
	SearchFacets map[string]any `json:"search_facets"` // Additional facets for filtering search results, e.g., by tags, groups, etc.
}

// ResourceSearch represents the result of a resource search query in CKAN.
type ResourceSearch struct {
	Count   int        `json:"count"`   // The total number of resources matching the search criteria
	Results []Resource `json:"results"` // The list of resources returned by the search
}

// TagSearch represents the result of a tag search query in CKAN.
type TagSearch struct {
	Count   int   `json:"count"`   // The total number of tags matching the search criteria
	Results []Tag `json:"results"` // The list of tags returned by the search
}
