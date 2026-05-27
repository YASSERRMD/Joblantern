package runner

import (
	"testing"
	"time"
)

func TestSummariseCountsMissed(t *testing.T) {
	r := Summarise(time.Now(), time.Now(), []Result{
		{FixtureID: "a", Detected: true, Source: "synthetic"},
		{FixtureID: "b", Detected: false, Source: "adversarial"},
		{FixtureID: "c", Detected: true, Source: "adversarial"},
	})
	if r.Total != 3 || r.Detected != 2 || r.Missed != 1 {
		t.Fatalf("unexpected: %+v", r)
	}
	if r.BySource["adversarial"] != 1 {
		t.Fatalf("adversarial detected count wrong: %d", r.BySource["adversarial"])
	}
}
