// Package groups organises operators into regional working groups
// so peers share lessons across instances.
package groups

// Region is a working group region.
type Region struct {
	ID       string
	Name     string
	Coordinator string
}

// Defaults returns the canonical regional working groups.
func Defaults() []Region {
	return []Region{
		{ID: "south-asia", Name: "South Asia (BD/IN/NP)", Coordinator: ""},
		{ID: "southeast-asia", Name: "Southeast Asia (PH/ID)", Coordinator: ""},
		{ID: "middle-east-corridor", Name: "Middle East Corridor (AE/SA/QA)", Coordinator: ""},
		{ID: "east-africa", Name: "East Africa", Coordinator: ""},
		{ID: "west-africa", Name: "West Africa", Coordinator: ""},
	}
}
