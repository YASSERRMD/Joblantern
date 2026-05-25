package risk_test

import (
	"testing"

	"github.com/yasserrmd/joblantern/internal/agent"
	"github.com/yasserrmd/joblantern/internal/risk"
)

func TestScore_Red(t *testing.T) {
	out := risk.Score([]agent.Fact{
		{SupportsRisk: "red", Weight: 0.9, FactType: "law.recruitment_fee_illegal"},
		{SupportsRisk: "red", Weight: 0.7, FactType: "pattern.red_flag"},
	}, risk.DefaultBands)
	if out.OverallRisk != "red" {
		t.Errorf("risk=%q", out.OverallRisk)
	}
}

func TestScore_Yellow(t *testing.T) {
	out := risk.Score([]agent.Fact{
		{SupportsRisk: "red", Weight: 0.5, FactType: "pattern.red_flag"},
		{SupportsRisk: "yellow", Weight: 0.2, FactType: "domain.freshness_score"},
	}, risk.DefaultBands)
	if out.OverallRisk != "yellow" {
		t.Errorf("risk=%q", out.OverallRisk)
	}
}

func TestScore_Green(t *testing.T) {
	out := risk.Score([]agent.Fact{
		{SupportsRisk: "green", Weight: 0.8, FactType: "registry.company_found"},
	}, risk.DefaultBands)
	if out.OverallRisk != "green" {
		t.Errorf("risk=%q", out.OverallRisk)
	}
}

func TestContradiction_DomainVsCompany(t *testing.T) {
	out := risk.Score([]agent.Fact{
		{FactType: "domain.age", Value: map[string]any{"age_days": 30.0}, SupportsRisk: "neutral"},
		{FactType: "registry.company_age", Value: map[string]any{"age_years": 10.0}, SupportsRisk: "neutral"},
	}, risk.DefaultBands)
	if len(out.Contradictions) != 1 {
		t.Fatalf("expected 1 contradiction, got %d", len(out.Contradictions))
	}
}
