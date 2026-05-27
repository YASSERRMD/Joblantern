// Package runner orchestrates the nightly red-team run. It pulls
// fixtures from generator + adversarial, sends them through the
// agent, and produces a detection-rate report.
package runner

import "time"

// Result is one fixture outcome.
type Result struct {
	FixtureID    string
	ExpectedBand string
	GotBand      string
	Detected     bool
	Latency      time.Duration
	Source       string // "synthetic" | "adversarial"
}

// Report is the nightly summary.
type Report struct {
	StartedAt time.Time
	Finished  time.Time
	Total     int
	Detected  int
	Missed    int
	BySource  map[string]int
}

// Summarise rolls results into a report.
func Summarise(start, end time.Time, results []Result) Report {
	r := Report{StartedAt: start, Finished: end, Total: len(results), BySource: map[string]int{}}
	for _, x := range results {
		if x.Detected {
			r.Detected++
			r.BySource[x.Source]++
		} else {
			r.Missed++
		}
	}
	return r
}
