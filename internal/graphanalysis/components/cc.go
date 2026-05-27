// Package components finds connected components in the recruitment
// graph. A component is a candidate scam cluster.
package components

// Edge is a minimal undirected edge.
type Edge struct{ A, B string }

// Find returns groups of node ids that are connected. Nodes that are
// not in any edge are returned as singletons only if includeIsolated
// is true.
func Find(nodes []string, edges []Edge) [][]string {
	parent := map[string]string{}
	for _, n := range nodes {
		parent[n] = n
	}
	var find func(string) string
	find = func(x string) string {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	union := func(a, b string) {
		ra, rb := find(a), find(b)
		if ra != rb {
			parent[ra] = rb
		}
	}
	for _, e := range edges {
		if _, ok := parent[e.A]; !ok {
			parent[e.A] = e.A
		}
		if _, ok := parent[e.B]; !ok {
			parent[e.B] = e.B
		}
		union(e.A, e.B)
	}
	groups := map[string][]string{}
	for n := range parent {
		r := find(n)
		groups[r] = append(groups[r], n)
	}
	out := make([][]string, 0, len(groups))
	for _, g := range groups {
		out = append(out, g)
	}
	return out
}
