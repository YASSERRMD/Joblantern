package agent_test

import (
	"context"
	"testing"

	"github.com/yasserrmd/joblantern/internal/agent"
)

type stubSub struct {
	name  string
	facts []agent.Fact
}

func (s *stubSub) Name() string                                           { return s.name }
func (s *stubSub) Run(_ context.Context, _ agent.Submission) []agent.Fact { return s.facts }

func TestOrchestrator_FanOut(t *testing.T) {
	subs := []agent.Subagent{
		&stubSub{name: "a", facts: []agent.Fact{{Source: "a", ToolName: "t1", FactType: "x", SupportsRisk: "red", Weight: 0.6}}},
		&stubSub{name: "b", facts: []agent.Fact{{Source: "b", ToolName: "t2", FactType: "y", SupportsRisk: "yellow", Weight: 0.5}}},
		&stubSub{name: "c", facts: []agent.Fact{{Source: "c", ToolName: "t3", FactType: "z", SupportsRisk: "green", Weight: 0.4}}},
	}
	orch := agent.New(subs...)
	v := orch.Run(context.Background(), agent.Submission{ID: "abc"})
	if len(v.Facts) != 3 {
		t.Fatalf("got %d facts", len(v.Facts))
	}
	if v.OverallRisk != "yellow" {
		t.Errorf("risk=%q want yellow", v.OverallRisk)
	}
	if v.Confidence <= 0 || v.Confidence > 1 {
		t.Errorf("confidence out of range: %v", v.Confidence)
	}
}

func TestOrchestrator_AllRed(t *testing.T) {
	subs := []agent.Subagent{
		&stubSub{name: "a", facts: []agent.Fact{{SupportsRisk: "red", Weight: 0.9}}},
		&stubSub{name: "b", facts: []agent.Fact{{SupportsRisk: "red", Weight: 0.9}}},
	}
	v := agent.New(subs...).Run(context.Background(), agent.Submission{})
	if v.OverallRisk != "red" {
		t.Errorf("risk=%q want red", v.OverallRisk)
	}
}
