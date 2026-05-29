// Package adversarial mutates known scam listings into evasive
// variants designed to slip past current rules.
package adversarial

import "strings"

// Mutator takes a base listing and produces a variant that aims to
// evade detection.
type Mutator interface {
	Mutate(body string) string
}

// PhoneSplitter splits a phone number across formatting so simple
// regex patterns miss it.
type PhoneSplitter struct{}

func (PhoneSplitter) Mutate(body string) string {
	return strings.ReplaceAll(body, "+971", "+ 9 7 1")
}

// HomographSwap replaces Latin letters with visually similar Unicode
// characters in suspicious vocabulary.
type HomographSwap struct{}

func (HomographSwap) Mutate(body string) string {
	r := strings.NewReplacer("URGENT", "URGEN"+string('Т'), "Pay", "Pаy")
	return r.Replace(body)
}

// WhitespaceInject inserts zero-width chars between trigger words.
type WhitespaceInject struct{}

func (WhitespaceInject) Mutate(body string) string {
	// Insert a zero-width space (U+200B) inside the trigger word so
	// naive substring matching misses it.
	return strings.ReplaceAll(body, "WhatsApp", "What\u200bsApp")
}

// All returns the canonical mutator chain.
func All() []Mutator { return []Mutator{PhoneSplitter{}, HomographSwap{}, WhitespaceInject{}} }
