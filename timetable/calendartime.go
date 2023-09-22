package timetable

import "time"

type CalendarTime struct {
	time.Time
}

func (c *CalendarTime) UnmarshalJSON(b []byte) error {
	t, err := time.ParseInLocation(`"2006-01-02T15:04:05"`, string(b), time.Local)
	if err != nil {
		return err
	}

	c.Time = t
	return nil
}

func (c *CalendarTime) MarshalJSON() ([]byte, error) {
	return []byte(c.Format(`"2006-01-02T15:04:05"`)), nil
}
