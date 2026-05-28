// Package provenance tracks the source of each triple. We use named
// graphs (the fourth element of a quad) so every fact retains its
// origin and harvest time.
package provenance

import "time"

// Quad is a triple plus the named-graph context.
type Quad struct {
	Subject     string
	Predicate   string
	Object      string
	Graph       string
	Source      string // e.g. "mcp-registry:gov-uk"
	HarvestedAt time.Time
}

// SourceIRI returns the canonical IRI for a source string.
func SourceIRI(source string) string { return "https://joblantern.org/kg/source/" + source }
