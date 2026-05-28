// Package page renders the optional public sustainability page at
// /sustainability. Numbers come from the footprint package.
package page

import "fmt"

// Data is what the page expects.
type Data struct {
	Month               string
	Verdicts            int
	TotalGramsCO2       float64
	PerVerdictGrams     float64
	GreenSchedulerHours int
	BatchInLowCarbonPct float64
}

// Render produces a deterministic HTML page.
func Render(d Data) string {
	return fmt.Sprintf(`<!doctype html>
<html><head><meta charset="utf-8"><title>Sustainability — Joblantern</title></head><body>
<h1>Sustainability — %s</h1>
<ul>
<li>Verdicts: %d</li>
<li>Total CO2-eq: %.0f g</li>
<li>Per verdict: %.2f g</li>
<li>Batch jobs run in low-carbon hours: %.0f%%</li>
</ul>
</body></html>`, d.Month, d.Verdicts, d.TotalGramsCO2, d.PerVerdictGrams, d.BatchInLowCarbonPct)
}
