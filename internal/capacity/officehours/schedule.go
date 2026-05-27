// Package officehours holds the recurring office-hours schedule.
package officehours

import "time"

// Session is one office-hours event.
type Session struct {
	Region   string
	Weekday  time.Weekday
	HourUTC  int
	Duration time.Duration
	URL      string
}

// Defaults returns the canonical recurring schedule.
func Defaults() []Session {
	return []Session{
		{Region: "south-asia", Weekday: time.Tuesday, HourUTC: 12, Duration: time.Hour},
		{Region: "southeast-asia", Weekday: time.Wednesday, HourUTC: 9, Duration: time.Hour},
		{Region: "middle-east-corridor", Weekday: time.Thursday, HourUTC: 13, Duration: time.Hour},
		{Region: "africa", Weekday: time.Monday, HourUTC: 15, Duration: time.Hour},
	}
}
