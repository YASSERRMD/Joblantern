// Package lib documents Joblantern's graph-library choice. Pulling in
// a large dependency like gonum/graph (Apache 2.0) is gated behind
// this small adapter so the bulk of the codebase remains free of
// vendor-specific types.
package lib

// Graph is the minimum interface the rest of the package programs
// against. Production wires this to gonum/graph/simple.WeightedGraph.
type Graph interface {
	AddNode(id string)
	AddEdge(from, to string, weight float64)
	Neighbors(id string) []string
	Weight(from, to string) (float64, bool)
	Nodes() []string
}
