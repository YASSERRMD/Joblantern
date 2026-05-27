// Package location flags improbable destination-vs-skill matches.
//
// Example: a welder offered a job in a country that does not import
// welders from the candidate's origin country is plausible; a
// software engineer offered relocation to that same country at the
// same salary is suspicious. The signal is the *combination* of
// destination, role, and origin.
package location

// Network captures the empirical destination-country networks for a
// role and an origin country. Populated from public migration data
// in deployment config.
type Network struct {
	Role          string
	OriginCountry string
	Destinations  []string // typical destinations
}

// Anomaly returns true if the destination is empirically unusual for
// this role + origin combination.
func Anomaly(role, origin, destination string, networks []Network) bool {
	for _, n := range networks {
		if n.Role != role || n.OriginCountry != origin {
			continue
		}
		for _, d := range n.Destinations {
			if d == destination {
				return false
			}
		}
		return true
	}
	return false // no data, do not flag
}
