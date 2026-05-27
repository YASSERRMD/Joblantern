// Package report is the server-side surface of the "report this
// number" post-call flow. The client app submits the hashed number,
// the user's quick category, and an optional free-text reason.
package report

import (
	"errors"
	"strings"
	"time"
)

// Category enumerates the quick-pick options shown in the app.
type Category string

const (
	CategoryScam       Category = "scam"
	CategoryAggressive Category = "aggressive"
	CategoryMistake    Category = "mistake"
	CategoryUnknown    Category = "unknown"
)

// Submission is the wire form.
type Submission struct {
	PhoneHash string    `json:"phone_hash"`
	Category  Category  `json:"category"`
	Reason    string    `json:"reason,omitempty"`
	At        time.Time `json:"at"`
	Country   string    `json:"country,omitempty"`
}

// Validate enforces minimum field requirements.
func (s Submission) Validate() error {
	if strings.TrimSpace(s.PhoneHash) == "" {
		return errors.New("phone hash required")
	}
	switch s.Category {
	case CategoryScam, CategoryAggressive, CategoryMistake, CategoryUnknown:
	default:
		return errors.New("unknown category")
	}
	return nil
}
