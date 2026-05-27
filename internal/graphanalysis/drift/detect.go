// Package drift detects when a graph component grows quickly,
// indicating a freshly active scam cluster.
package drift

import "time"

// Snapshot is the per-component count at a moment in time.
type Snapshot struct {
	ComponentID string
	At          time.Time
	NodeCount   int
}

// Alert is emitted when a component grows more than a multiplicative
// factor across a window.
type Alert struct {
	ComponentID string
	From        Snapshot
	To          Snapshot
	GrowthRatio float64
}

// Detect flags components whose count more than multiplier-fold in the
// last window. Snapshots are assumed sorted oldest-first per component.
func Detect(s []Snapshot, window time.Duration, multiplier float64) []Alert {
	by := map[string][]Snapshot{}
	for _, x := range s {
		by[x.ComponentID] = append(by[x.ComponentID], x)
	}
	var out []Alert
	for id, xs := range by {
		if len(xs) < 2 {
			continue
		}
		latest := xs[len(xs)-1]
		var ref *Snapshot
		for i := len(xs) - 2; i >= 0; i-- {
			if latest.At.Sub(xs[i].At) >= window {
				ref = &xs[i]
				break
			}
		}
		if ref == nil {
			continue
		}
		if ref.NodeCount == 0 {
			continue
		}
		ratio := float64(latest.NodeCount) / float64(ref.NodeCount)
		if ratio >= multiplier {
			out = append(out, Alert{ComponentID: id, From: *ref, To: latest, GrowthRatio: ratio})
		}
	}
	return out
}
