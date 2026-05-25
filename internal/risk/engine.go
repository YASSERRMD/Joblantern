// Package risk is the deterministic risk-scoring engine. The agent
// gathers evidence and the LLM proposes a narrative; this engine
// computes the numeric verdict so the same input always yields the
// same output.
package risk

import (
	"sort"

	"github.com/yasserrmd/joblantern/internal/agent"
)

// Bands are tunable thresholds that move risk between green / yellow /
// red. They sum weighted red and yellow facts and compare against
// these thresholds.
type Bands struct {
	RedScore    float64 // total weighted-red ≥ this → red
	YellowScore float64 // total weighted (red+yellow) ≥ this → yellow
}

// DefaultBands match the placeholder in agent.score; the risk engine
// also adds contradiction detection and richer reason text.
var DefaultBands = Bands{RedScore: 0.9, YellowScore: 0.6}

// Output is the engine's return value. The agent stores it on the
// Verdict in place of the placeholder fields it populated earlier.
type Output struct {
	OverallRisk    string             `json:"overall_risk"`
	Confidence     float64            `json:"confidence"`
	Reasons        []string           `json:"reasons"`
	Contradictions []Contradiction    `json:"contradictions,omitempty"`
	WeightApplied  map[string]float64 `json:"weight_applied"`
}

// Contradiction is a declarative inconsistency between two facts.
type Contradiction struct {
	A      string `json:"a"`
	B      string `json:"b"`
	Reason string `json:"reason"`
}

// Score the facts deterministically. Sort facts before iterating so
// JSON output is reproducible.
func Score(facts []agent.Fact, bands Bands) Output {
	sorted := make([]agent.Fact, len(facts))
	copy(sorted, facts)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Source != sorted[j].Source {
			return sorted[i].Source < sorted[j].Source
		}
		if sorted[i].ToolName != sorted[j].ToolName {
			return sorted[i].ToolName < sorted[j].ToolName
		}
		return sorted[i].FactType < sorted[j].FactType
	})

	var red, yellow, green float64
	weights := make(map[string]float64, len(sorted))
	var reasons []string
	for _, f := range sorted {
		weights[f.Source+"."+f.FactType] += f.Weight
		switch f.SupportsRisk {
		case "red":
			red += f.Weight
			reasons = append(reasons, f.FactType+" (red)")
		case "yellow":
			yellow += f.Weight
			reasons = append(reasons, f.FactType+" (yellow)")
		case "green":
			green += f.Weight
		}
	}
	risk := "green"
	switch {
	case red >= bands.RedScore:
		risk = "red"
	case red >= bands.YellowScore-0.2 || red+yellow >= bands.YellowScore:
		risk = "yellow"
	}
	conf := 0.2
	if n := len(sorted); n > 0 {
		conf = 0.2 + float64(n)*0.05
		if conf > 0.95 {
			conf = 0.95
		}
	}
	contradictions := detectContradictions(sorted)
	return Output{
		OverallRisk:    risk,
		Confidence:     conf,
		Reasons:        reasons,
		Contradictions: contradictions,
		WeightApplied:  weights,
	}
}

// detectContradictions runs the declared inconsistency rules over the
// fact set. Each rule is conservative — it only fires when both sides
// of the contradiction are explicitly present.
func detectContradictions(facts []agent.Fact) []Contradiction {
	var out []Contradiction
	var (
		domainAgeDays int
		hasDomainAge  bool
		companyAgeYrs int
		hasCompanyAge bool
	)
	for _, f := range facts {
		if f.FactType == "domain.age" {
			if m, ok := f.Value.(map[string]any); ok {
				if v, ok := m["age_days"].(float64); ok {
					domainAgeDays = int(v)
					hasDomainAge = true
				}
			}
		}
		if f.FactType == "registry.company_age" {
			if m, ok := f.Value.(map[string]any); ok {
				if v, ok := m["age_years"].(float64); ok {
					companyAgeYrs = int(v)
					hasCompanyAge = true
				}
			}
		}
	}
	if hasDomainAge && hasCompanyAge && domainAgeDays < 90 && companyAgeYrs > 5 {
		out = append(out, Contradiction{
			A: "domain.age", B: "registry.company_age",
			Reason: "Domain registered <90 days ago but company claims >5 years operation.",
		})
	}
	return out
}
