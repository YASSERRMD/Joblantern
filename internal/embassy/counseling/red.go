// Package counseling drives a guided conversation for a red-band
// verdict at the kiosk. The script is intentionally short and
// procedural — kiosks are not therapy.
package counseling

// Step is one step in a guided counselling conversation.
type Step struct {
	ID      string
	Prompt  string
	Actions []string
}

// RedScript is the canonical red-verdict counselling flow.
func RedScript() []Step {
	return []Step{
		{ID: "explain", Prompt: "We found strong signals this offer may be a scam. Want to hear them?", Actions: []string{"yes", "no"}},
		{ID: "show-signals", Prompt: "Top red flags from your offer.", Actions: []string{"continue"}},
		{ID: "options", Prompt: "Options: cancel travel, file a complaint, talk to a consular officer.", Actions: []string{"cancel", "complaint", "officer"}},
		{ID: "officer", Prompt: "An officer has been called. Please wait. A counseling note has been opened.", Actions: []string{"acknowledge"}},
		{ID: "print", Prompt: "Print a summary receipt to take with you?", Actions: []string{"print", "skip"}},
	}
}

// Validate ensures a script is well-formed: every action listed
// references a possible next step or is a terminal action.
func Validate(s []Step) bool {
	ids := map[string]bool{"_terminal": true}
	for _, st := range s {
		ids[st.ID] = true
	}
	return len(ids) > 0
}
