// Package advisory is the research advisory board. The board is
// separate from the Trust & Safety council so research and safety
// concerns can be balanced rather than collapsed.
package advisory

// Member is one board member.
type Member struct {
	ID          string
	Name        string
	Field       string // "migration", "labor-econ", "law", "ml-ethics", "data-engineering"
	Affiliation string
	Term        string // "2026-2028"
}

// Defaults captures the seed composition.
func Defaults() []Member {
	return []Member{
		{ID: "m1", Field: "migration"},
		{ID: "m2", Field: "labor-econ"},
		{ID: "m3", Field: "law"},
		{ID: "m4", Field: "ml-ethics"},
		{ID: "m5", Field: "data-engineering"},
	}
}
