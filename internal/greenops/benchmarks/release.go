// Package benchmarks holds per-release efficiency benchmarks. Each
// release adds a row; regressions are blocking.
package benchmarks

import "time"

// Row is one per-release measurement.
type Row struct {
	Release            string
	CapturedAt         time.Time
	JoulesPerVerdict   float64
	TokensPerVerdict   int
	MCPCallsPerVerdict int
	GramsCO2PerVerdict float64
}

// Regression is the per-axis tolerance.
type Regression struct {
	JoulesPctMax   float64
	TokensPctMax   float64
	MCPCallsPctMax float64
	CO2PctMax      float64
}

// Default is the tolerance Joblantern publishes — small regressions
// are allowed if explained in the release notes; large regressions
// fail CI.
func Default() Regression {
	return Regression{JoulesPctMax: 10, TokensPctMax: 10, MCPCallsPctMax: 5, CO2PctMax: 10}
}
