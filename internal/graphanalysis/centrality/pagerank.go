// Package centrality implements a basic damped-PageRank scoring so
// the investigator workbench can highlight the ringleader nodes of
// each component.
package centrality

// Edge is a directed edge with weight.
type Edge struct {
	From   string
	To     string
	Weight float64
}

// PageRank returns rank scores for each node id. Damping defaults to
// 0.85 when 0. Iter defaults to 30.
func PageRank(nodes []string, edges []Edge, damping float64, iter int) map[string]float64 {
	if damping <= 0 {
		damping = 0.85
	}
	if iter <= 0 {
		iter = 30
	}
	N := float64(len(nodes))
	if N == 0 {
		return map[string]float64{}
	}
	out := map[string]float64{}
	for _, n := range nodes {
		out[n] = 1 / N
	}
	outAdj := map[string][]Edge{}
	outSum := map[string]float64{}
	for _, e := range edges {
		outAdj[e.From] = append(outAdj[e.From], e)
		outSum[e.From] += e.Weight
	}
	for i := 0; i < iter; i++ {
		next := map[string]float64{}
		for _, n := range nodes {
			next[n] = (1 - damping) / N
		}
		for from, es := range outAdj {
			share := damping * out[from] / outSum[from]
			for _, e := range es {
				next[e.To] += share * e.Weight
			}
		}
		out = next
	}
	return out
}
