// Package agent is the Joblantern orchestrator. It accepts a job
// submission, fans out to MCP-backed sub-agents in parallel,
// aggregates evidence, and emits a Verdict.
//
// Note on ADK Go: Phase 13 ships the orchestration logic in pure-Go
// using goroutines + channels. A future change can swap the internal
// scheduler for google.golang.org/adk's LLMAgent / ParallelAgent
// without changing the public types or the HTTP surface.
package agent

import (
	"context"
	"time"
)

// Submission is the incoming verification request.
type Submission struct {
	ID             string  `json:"id,omitempty"`
	UserID         string  `json:"user_id,omitempty"`
	ListingURL     string  `json:"listing_url,omitempty"`
	ListingText    string  `json:"listing_text,omitempty"`
	CompanyName    string  `json:"company_name,omitempty"`
	ClaimedAddress string  `json:"claimed_address,omitempty"`
	RecruiterEmail string  `json:"recruiter_email,omitempty"`
	RecruiterPhone string  `json:"recruiter_phone,omitempty"`
	Role           string  `json:"role,omitempty"`
	ClaimedSalary  float64 `json:"claimed_salary,omitempty"`
	SalaryCurrency string  `json:"salary_currency,omitempty"`
	SalaryPeriod   string  `json:"salary_period,omitempty"`
	Jurisdiction   string  `json:"jurisdiction,omitempty"`
	Domain         string  `json:"domain,omitempty"`
	HomeLat        float64 `json:"home_lat,omitempty"`
	HomeLon        float64 `json:"home_lon,omitempty"`
}

// Fact is one piece of evidence the agent gathered.
type Fact struct {
	Source       string  `json:"source"`
	ToolName     string  `json:"tool_name"`
	FactType     string  `json:"fact_type"`
	Value        any     `json:"value"`
	SupportsRisk string  `json:"supports_risk"`
	Weight       float64 `json:"weight"`
	Citation     string  `json:"citation,omitempty"`
}

// Verdict is the agent's final structured output.
type Verdict struct {
	VerificationID string    `json:"verification_id"`
	OverallRisk    string    `json:"overall_risk"`
	Confidence     float64   `json:"confidence"`
	Facts          []Fact    `json:"facts"`
	Reasons        []string  `json:"reasons"`
	GeneratedAt    time.Time `json:"generated_at"`
}

// Subagent is implemented by each topical investigator. Each is run
// concurrently by the orchestrator.
type Subagent interface {
	Name() string
	Run(ctx context.Context, sub Submission) []Fact
}
