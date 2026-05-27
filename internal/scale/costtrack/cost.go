// Package costtrack records per-verdict cost so we can publish a
// cost-per-verdict dashboard. Costs include LLM inference, MCP
// queries, storage, and bandwidth.
package costtrack

import "time"

// Sample is one measurement.
type Sample struct {
	VerdictID    string
	At           time.Time
	LLMCents     float64
	MCPCents     float64
	StorageCents float64
	BandwidthCents float64
}

// Total returns the total cents-cost for this verdict.
func (s Sample) Total() float64 {
	return s.LLMCents + s.MCPCents + s.StorageCents + s.BandwidthCents
}

// MonthlyRollup is the aggregate.
type MonthlyRollup struct {
	Month       time.Month
	Year        int
	VerdictCount int
	TotalCents  float64
	MedianCents float64
}
