// Package scheduler defers low-priority batch jobs (red-team runs,
// aggregations, model training) to low-carbon hours.
package scheduler

import "time"

// Window is one scheduled window during the next 24h.
type Window struct {
	Start time.Time
	End   time.Time
	Score float64 // lower is greener
}

// Pick returns the lowest-carbon window that fits a job of `duration`
// from the supplied forecast. If no window fits, returns now.
func Pick(now time.Time, duration time.Duration, candidates []Window) Window {
	if duration <= 0 || len(candidates) == 0 {
		return Window{Start: now, End: now.Add(duration)}
	}
	var best Window
	bestScore := -1.0
	for _, w := range candidates {
		if w.End.Sub(w.Start) < duration {
			continue
		}
		if bestScore < 0 || w.Score < bestScore {
			best = w
			bestScore = w.Score
		}
	}
	if bestScore < 0 {
		return Window{Start: now, End: now.Add(duration)}
	}
	return best
}
