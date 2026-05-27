// Package colocation picks colocation regions by carbon intensity
// when the data-residency policy gives us a choice.
package colocation

// Region is a colocation candidate.
type Region struct {
	ID            string
	Name          string
	GridGCO2PerKWh float64
}

// Greenest returns the lowest-carbon region from a permissible set.
func Greenest(allowed []Region) Region {
	if len(allowed) == 0 {
		return Region{}
	}
	best := allowed[0]
	for _, r := range allowed[1:] {
		if r.GridGCO2PerKWh < best.GridGCO2PerKWh {
			best = r
		}
	}
	return best
}
