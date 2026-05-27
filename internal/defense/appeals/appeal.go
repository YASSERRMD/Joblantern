// Package appeals is the false-positive remedy for an entity that
// was wrongly listed. Appeals are reviewed by the Trust & Safety
// Council and, if upheld, trigger a reversal across all downstream
// surfaces (blocklist, safebrowsing, regulator forwards).
package appeals

import (
	"errors"
	"strings"
	"time"
)

// Status is the lifecycle.
type Status string

const (
	StatusOpen      Status = "open"
	StatusUpheld    Status = "upheld"
	StatusRejected  Status = "rejected"
	StatusWithdrawn Status = "withdrawn"
)

// Appeal is one filing.
type Appeal struct {
	ID            string
	Subject       string
	FiledBy       string
	FiledAt       time.Time
	ReceivedBands []string
	Argument      string
	Status        Status
	DecidedAt     time.Time
	DecisionNote  string
}

// Validate checks the appeal has the minimum content.
func (a Appeal) Validate() error {
	if strings.TrimSpace(a.Subject) == "" {
		return errors.New("subject required")
	}
	if len([]rune(strings.TrimSpace(a.Argument))) < 100 {
		return errors.New("argument must be at least 100 characters")
	}
	return nil
}
