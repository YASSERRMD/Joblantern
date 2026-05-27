// Package patterns matches new listings against confirmed-scam
// domains and their close variants.
package patterns

import "strings"

// Confirmed is the maintained set of confirmed-scam domains.
type Confirmed map[string]struct{}

// Match returns true if the listing's host equals a confirmed domain
// or differs by a single-character edit (typo squat).
func (c Confirmed) Match(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	if _, ok := c[host]; ok {
		return true
	}
	for d := range c {
		if editDistance(host, d) == 1 {
			return true
		}
	}
	return false
}

func editDistance(a, b string) int {
	if a == b {
		return 0
	}
	la, lb := len(a), len(b)
	if abs(la-lb) > 1 {
		return 2 // we only care about distance ∈ {0,1}, anything else is "far"
	}
	i, j := 0, 0
	diff := 0
	for i < la && j < lb {
		if a[i] == b[j] {
			i++
			j++
			continue
		}
		diff++
		if diff > 1 {
			return diff
		}
		if la == lb {
			i++
			j++
		} else if la > lb {
			i++
		} else {
			j++
		}
	}
	if i < la || j < lb {
		diff++
	}
	return diff
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
