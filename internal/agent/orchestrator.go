package agent

import (
	"context"
	"sort"
	"sync"
	"time"
)

// ScoreFunc is the pluggable scoring engine. Replaced in main.go with
// the real risk-engine scorer (internal/risk.Score).
type ScoreFunc func(facts []Fact) (risk string, confidence float64, reasons []string)

// Orchestrator runs all configured sub-agents in parallel and
// aggregates the facts they emit. Final scoring is performed by the
// internal/risk engine (Phase 14); for Phase 13 the orchestrator
// emits a verdict with a simple aggregation so the HTTP surface has
// something deterministic to return.
type Orchestrator struct {
	Subagents []Subagent
	Timeout   time.Duration
	Score     ScoreFunc // optional; defaults to the built-in placeholder
}

// New returns an Orchestrator with sensible defaults.
func New(subs ...Subagent) *Orchestrator {
	return &Orchestrator{Subagents: subs, Timeout: 60 * time.Second}
}

// WithScorer overrides the default placeholder scorer.
func (o *Orchestrator) WithScorer(f ScoreFunc) *Orchestrator {
	o.Score = f
	return o
}

// Run executes every sub-agent concurrently and returns a Verdict.
func (o *Orchestrator) Run(ctx context.Context, sub Submission) Verdict {
	ctx, cancel := context.WithTimeout(ctx, o.Timeout)
	defer cancel()

	resCh := make(chan []Fact, len(o.Subagents))
	var wg sync.WaitGroup
	for _, sa := range o.Subagents {
		wg.Add(1)
		go func(s Subagent) {
			defer wg.Done()
			resCh <- s.Run(ctx, sub)
		}(sa)
	}
	wg.Wait()
	close(resCh)

	var all []Fact
	for facts := range resCh {
		all = append(all, facts...)
	}
	sort.SliceStable(all, func(i, j int) bool {
		if all[i].Source != all[j].Source {
			return all[i].Source < all[j].Source
		}
		return all[i].ToolName < all[j].ToolName
	})

	scorer := o.Score
	if scorer == nil {
		scorer = score
	}
	risk, conf, reasons := scorer(all)
	return Verdict{
		VerificationID: sub.ID,
		OverallRisk:    risk,
		Confidence:     conf,
		Facts:          all,
		Reasons:        reasons,
		GeneratedAt:    time.Now().UTC(),
	}
}

// score is a placeholder; the real risk engine in Phase 14 replaces it.
func score(facts []Fact) (string, float64, []string) {
	var red, yellow, green float64
	var reasons []string
	for _, f := range facts {
		switch f.SupportsRisk {
		case "red":
			red += f.Weight
			reasons = append(reasons, f.FactType+" (red)")
		case "yellow":
			yellow += f.Weight
		case "green":
			green += f.Weight
		}
	}
	risk := "green"
	switch {
	case red >= 0.9:
		risk = "red"
	case red >= 0.4 || red+yellow >= 0.6:
		risk = "yellow"
	}
	conf := 0.2
	if len(facts) > 0 {
		conf = 0.2 + float64(len(facts))*0.05
		if conf > 0.95 {
			conf = 0.95
		}
	}
	return risk, conf, reasons
}
