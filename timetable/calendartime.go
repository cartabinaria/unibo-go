// SPDX-FileCopyrightText: 2023 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: MIT

package timetable

import (
	"fmt"
	"time"
)

// CalendarTime is a time.Time that has a custom MarshalJSON and UnmarshalJSON.
//
// It is used to parse the timetable, since the endpoint uses a custom format for dates.
//
// It always uses the Italian timezone.
type CalendarTime struct {
	time.Time
}

const (
	layout       = `"2006-01-02T15:04:05"` // The layout of the dates in the timetable as returned by the API
	ianaTimezone = "Europe/Rome"           // The IANA timezone of the dates in the timetable
)

var cachedTimezone *time.Location

func cacheTimezone() (err error) {
	if cachedTimezone != nil {
		return
	}

	timezone, err := time.LoadLocation(ianaTimezone)
	if err != nil {
		return
	}

	cachedTimezone = timezone
	return
}

func (c *CalendarTime) UnmarshalJSON(b []byte) error {
	err := cacheTimezone()
	if err != nil {
		return fmt.Errorf("could not load italian timezone: %w", err)
	}

	t, err := time.ParseInLocation(layout, string(b), cachedTimezone)
	if err != nil {
		return err
	}

	c.Time = t
	return nil
}

func (c *CalendarTime) MarshalJSON() ([]byte, error) {
	return []byte(c.Format(layout)), nil
}
