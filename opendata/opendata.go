// SPDX-FileCopyrightText: 2023 - 2025 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: MIT

// Package opendata provides functions to fetch data from the UniBo Open Data
// portal.
package opendata

import "github.com/cartabinaria/unibo-go/ckan"

const openDataUrl = "https://dati.unibo.it"

var ckanClient = ckan.NewClient(openDataUrl)
