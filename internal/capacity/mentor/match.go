// Package mentor pairs new operators with experienced ones based on
// shared language and region.
package mentor

import "strings"

// Operator is one operator profile.
type Operator struct {
	ID        string
	Country   string
	Region    string
	Languages []string
	IsMentor  bool
	OpenSlots int
}

// Match returns a mentor for the new operator, or empty string if none
// fits. Strategy: same region first, then shared language.
func Match(newOp Operator, pool []Operator) string {
	var sameRegion, sharedLang []Operator
	for _, p := range pool {
		if !p.IsMentor || p.OpenSlots <= 0 {
			continue
		}
		if strings.EqualFold(p.Region, newOp.Region) {
			sameRegion = append(sameRegion, p)
		}
		for _, l := range p.Languages {
			for _, m := range newOp.Languages {
				if strings.EqualFold(l, m) {
					sharedLang = append(sharedLang, p)
				}
			}
		}
	}
	if len(sameRegion) > 0 {
		return sameRegion[0].ID
	}
	if len(sharedLang) > 0 {
		return sharedLang[0].ID
	}
	return ""
}
