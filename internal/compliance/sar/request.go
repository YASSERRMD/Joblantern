// Package sar processes Subject Access Requests: export, correction,
// and deletion. The SLA we commit to is 30 days from intake.
package sar

import (
	"errors"
	"time"
)

// Kind enumerates the SAR operations.
type Kind string

const (
	KindExport   Kind = "export"
	KindCorrect  Kind = "correct"
	KindDelete   Kind = "delete"
	KindObject   Kind = "object" // object to processing
	KindRestrict Kind = "restrict"
)

// Request is one SAR.
type Request struct {
	ID          string
	Subject     string
	Kind        Kind
	FiledAt     time.Time
	DueAt       time.Time
	VerifiedAt  time.Time
	CompletedAt time.Time
	Notes       string
}

// SLA is the response window.
const SLA = 30 * 24 * time.Hour

// Compute returns the request with DueAt and basic validation.
func Compute(r Request) (Request, error) {
	if r.Subject == "" {
		return r, errors.New("subject required")
	}
	if r.FiledAt.IsZero() {
		r.FiledAt = time.Now().UTC()
	}
	r.DueAt = r.FiledAt.Add(SLA)
	return r, nil
}
