// Package webhook delivers high-confidence scam verdicts to opted-in
// researcher endpoints. Deliveries are HMAC-signed and retried with
// exponential backoff.
package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Subscription describes a researcher's webhook target.
type Subscription struct {
	ID                  string
	Owner               string
	URL                 string
	Secret              []byte
	MinConfidence       float64
	IndustryFilter      []string
	CountryFilter       []string
	CreatedAt           time.Time
	LastSuccessfulPush  time.Time
	ConsecutiveFailures int
}

// Sign returns the HMAC-SHA256 hex digest of the payload using the
// subscription secret. Receivers verify by computing the same digest.
func (s Subscription) Sign(payload []byte) string {
	m := hmac.New(sha256.New, s.Secret)
	_, _ = m.Write(payload)
	return hex.EncodeToString(m.Sum(nil))
}

// Backoff returns the delay before the n-th retry attempt.
func Backoff(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	if attempt > 8 {
		attempt = 8
	}
	return time.Duration(1<<attempt) * time.Second
}
