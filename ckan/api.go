// SPDX-FileCopyrightText: 2025 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: MIT

package ckan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	httpClient = &http.Client{} // The HTTP client used to make requests
)

// ApiError represents an error returned by the CKAN API.
type ApiError struct {
	Message string `json:"message"` // The error message returned by the API
	Type    string `json:"__type"`  // The type of error, e.g., "NotFound", "ValidationError"
}

// ApiResponse is a generic structure for API responses from the CKAN API.
//
// Every API call returns an ApiResponse, which can either contain a successful result
// or an error. If the Success field is false, the Result field will be nil and the
// Error field will contain details about the failure.
//
// For a more idiomatic Go experience, consider using the Request function,
// which returns the result and error as separate return values.
type ApiResponse[T any] struct {
	Help    string    `json:"help"`             // The documentation URL for the called API endpoint
	Success bool      `json:"success"`          // Whether the API call was successful.
	Result  *T        `json:"result,omitempty"` // The result of the API call, which can be any type
	Error   *ApiError `json:"error,omitempty"`  // The error information if the API call failed
}

// RequestRaw queries the CKAN API and returns the raw ApiResponse.
//
// For a more idiomatic Go experience, consider using Request.
func RequestRaw[T any](url string) (*ApiResponse[T], error) {
	res, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to make request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var response ApiResponse[T]
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("unable to decode response: %w", err)
	}

	err = res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to close response body: %w", err)
	}

	return &response, nil
}

// Request queries the CKAN API and returns the result of type T.
//
// Please use the Client methods to interact with the CKAN API if
// you are not looking for a raw response.
//
//	client := ckan.NewClient("https://demo.ckan.org")
//	pkgs, err := client.GetPackageList()
//	...
//
// If the API call was not successful, it returns a standard error
// containing the error message from the API response.
//
// You should not need to use this function directly in most cases,
// as it is primarily used internally by the Client methods.
//
// For a more advanced use case, consider using RequestRaw.
func Request[T any](url string) (*T, error) {
	resp, err := RequestRaw[T](url)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		if resp.Error == nil {
			return nil, fmt.Errorf("API call was not successful, but no error information provided")
		}
		return nil, fmt.Errorf("API call failed: %s", resp.Error.Message)
	}

	return resp.Result, nil
}
