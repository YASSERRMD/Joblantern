// Package viz emits a d3-force compatible JSON layout.
package viz

import "encoding/json"

// Node is the d3 representation.
type Node struct {
	ID    string  `json:"id"`
	Group string  `json:"group,omitempty"`
	Rank  float64 `json:"rank,omitempty"`
}

// Link is the d3 representation.
type Link struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Value  float64 `json:"value"`
}

// Document is the full payload returned by the viz endpoint.
type Document struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

// Render produces a stable, deterministic JSON byte slice.
func Render(d Document) ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}
