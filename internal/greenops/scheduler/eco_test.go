package scheduler

import (
	"testing"
	"time"
)

func TestPicksLowestCarbonWindow(t *testing.T) {
	now := time.Date(2026, 5, 27, 8, 0, 0, 0, time.UTC)
	windows := []Window{
		{Start: now.Add(1 * time.Hour), End: now.Add(3 * time.Hour), Score: 350},
		{Start: now.Add(4 * time.Hour), End: now.Add(7 * time.Hour), Score: 120}, // greenest
		{Start: now.Add(10 * time.Hour), End: now.Add(12 * time.Hour), Score: 220},
	}
	got := Pick(now, 2*time.Hour, windows)
	if got.Score != 120 {
		t.Fatalf("expected score 120, got %v", got.Score)
	}
}

func TestFallsBackToNowWhenNoFit(t *testing.T) {
	now := time.Now()
	got := Pick(now, 10*time.Hour, []Window{{Start: now, End: now.Add(2 * time.Hour), Score: 100}})
	if got.Start != now {
		t.Fatalf("expected fallback to now")
	}
}
