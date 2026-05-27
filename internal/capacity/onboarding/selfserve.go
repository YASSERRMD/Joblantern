// Package onboarding orchestrates the self-service partner flow.
//
// The shipped target: a new NGO in Nepal completes onboarding,
// deploys their instance, and reaches first verdict within one
// working day.
package onboarding

import "time"

// Step is one step in the partner self-serve flow.
type Step struct {
	ID          string
	Title       string
	EstMinutes  int
	Order       int
}

// Defaults returns the canonical ordered steps.
func Defaults() []Step {
	return []Step{
		{ID: "register", Title: "Register the partner organisation", Order: 1, EstMinutes: 5},
		{ID: "verify-email", Title: "Verify partner contact email", Order: 2, EstMinutes: 2},
		{ID: "kyc-lite", Title: "Submit lite-KYC documents", Order: 3, EstMinutes: 10},
		{ID: "deploy", Title: "Deploy your instance (deploy-in-a-box)", Order: 4, EstMinutes: 30},
		{ID: "first-verdict", Title: "Run your first verdict", Order: 5, EstMinutes: 10},
		{ID: "mentor", Title: "Get matched with a mentor", Order: 6, EstMinutes: 10},
	}
}

// TotalEstimate returns the total estimated time-to-first-verdict.
func TotalEstimate() time.Duration {
	mins := 0
	for _, s := range Defaults() {
		mins += s.EstMinutes
	}
	return time.Duration(mins) * time.Minute
}
