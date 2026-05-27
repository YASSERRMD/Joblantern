// Package opendata describes the open-data release cadence.
package opendata

import "time"

// Release is one scheduled open-data release.
type Release struct {
	Cadence    string
	Window     time.Duration
	Format     string
}

// Defaults captures the schedule.
func Defaults() []Release {
	return []Release{
		{Cadence: "monthly", Window: 30 * 24 * time.Hour, Format: "parquet"},
		{Cadence: "annual", Window: 365 * 24 * time.Hour, Format: "parquet+json-ld"},
	}
}
