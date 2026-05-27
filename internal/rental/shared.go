// Package rental is the housing-scam vertical. Most of the agent
// infrastructure is shared with the job-offer vertical — the same
// MCP servers (registry, domain, salary, routing, scam-db) score
// rental contexts with a different rule pack.
package rental

import "context"

// Verdict is a housing-specific verdict.
type Verdict struct {
	ID         string
	RiskScore  int
	RiskBand   string
	RedFlags   []string
	Citations  []string
	JoinedJobID string
}

// Submission is the input to the rental agent.
type Submission struct {
	Country        string
	City           string
	ListingURL     string
	ListingText    string
	MonthlyRent    float64
	Currency       string
	DepositMethod  string
	ContactPhone   string
	ContactEmail   string
	LandlordName   string
	ImageURLs      []string
}

// AgentRunner is the contract the rental agent expects from the
// shared MCP-backed agent orchestrator.
type AgentRunner interface {
	Run(ctx context.Context, kind string, payload any) (any, error)
}
