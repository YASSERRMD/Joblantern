// Package funding enforces the no-single-funder-controls-council
// rule. If any one funder ever exceeds 25% of unrestricted income,
// the council is notified and the next budget shifts to diversify.
package funding

// Funder is one income source.
type Funder struct {
	ID     string
	Name   string
	Annual float64
}

// IndependenceThreshold is the maximum allowed share for a single
// funder of unrestricted income.
const IndependenceThreshold = 0.25

// Check returns the funder(s) exceeding the threshold.
func Check(funders []Funder) []Funder {
	total := 0.0
	for _, f := range funders {
		total += f.Annual
	}
	if total == 0 {
		return nil
	}
	var over []Funder
	for _, f := range funders {
		if f.Annual/total > IndependenceThreshold {
			over = append(over, f)
		}
	}
	return over
}
