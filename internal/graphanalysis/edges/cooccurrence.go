// Package edges builds edge weights from entity co-occurrence in
// scam reports. The more reports two entities appear in together,
// the heavier the edge.
package edges

// Edge is one weighted graph edge.
type Edge struct {
	From   string
	To     string
	Weight float64
	Count  int
}

// Builder accumulates co-occurrence counts.
type Builder struct {
	counts map[[2]string]int
}

// New constructs a fresh Builder.
func New() *Builder { return &Builder{counts: map[[2]string]int{}} }

// Observe records a co-occurrence between a and b.
func (b *Builder) Observe(a, c string) {
	if a == c {
		return
	}
	if a > c {
		a, c = c, a
	}
	b.counts[[2]string{a, c}]++
}

// Edges returns all observed edges. Weight is log(count+1) so
// runaway hubs do not dominate the layout.
func (b *Builder) Edges() []Edge {
	out := make([]Edge, 0, len(b.counts))
	for k, c := range b.counts {
		out = append(out, Edge{From: k[0], To: k[1], Count: c, Weight: log1p(float64(c))})
	}
	return out
}

func log1p(x float64) float64 {
	if x < 0 {
		return 0
	}
	// Avoid pulling math; small approximation suffices for ranking.
	switch {
	case x < 1:
		return x
	case x < 5:
		return 1 + (x-1)*0.5
	default:
		return 2 + (x-5)*0.1
	}
}
