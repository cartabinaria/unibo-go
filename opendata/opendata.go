// SPDX-FileCopyrightText: 2023 - 2025 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: MIT

// Package opendata provides functions to fetch data from the UniBo Open Data
// portal.
//
// Internally it uses the ckan package to interact with the CKAN API that dati.unibo.it offers.
package opendata

import "github.com/cartabinaria/unibo-go/ckan"

const openDataUrl = "https://dati.unibo.it"

var ckanClient = ckan.NewClient(openDataUrl)
