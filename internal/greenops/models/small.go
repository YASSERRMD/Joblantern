// Package models encourages compact LLM providers. Every verdict
// stage names the smallest model that does the job; bigger models
// are reserved for tie-breakers.
package models

// Choice is the policy for one stage.
type Choice struct {
	Stage   string
	Primary string
	Tiebreak string
}

// Defaults captures the policy.
func Defaults() []Choice {
	return []Choice{
		{Stage: "extract-redflags", Primary: "small-7b-quantized", Tiebreak: "medium-13b"},
		{Stage: "paraphrase-evidence", Primary: "small-7b-quantized"},
		{Stage: "explain-personalization", Primary: "small-7b-quantized"},
		{Stage: "complex-cross-link-reasoning", Primary: "medium-13b", Tiebreak: "large-70b"},
	}
}
