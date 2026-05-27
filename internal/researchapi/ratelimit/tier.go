// Package ratelimit enforces per-tier request budgets for the
// research API. Implementations use a token bucket per token subject.
package ratelimit

import (
	"sync"
	"time"
)

// Tier limits.
const (
	publicRPM     = 30
	academicRPM   = 600
	journalistRPM = 300
	regulatorRPM  = 1200
)

// LimitFor returns requests-per-minute for the given tier.
func LimitFor(tier string) int {
	switch tier {
	case "academic":
		return academicRPM
	case "journalist":
		return journalistRPM
	case "regulator":
		return regulatorRPM
	default:
		return publicRPM
	}
}

// Bucket is a thread-safe token bucket.
type Bucket struct {
	mu         sync.Mutex
	tokens     float64
	maxTokens  float64
	refillRate float64
	last       time.Time
}

// NewBucket creates a bucket with the supplied per-minute rate.
func NewBucket(rpm int) *Bucket {
	rate := float64(rpm) / 60.0
	return &Bucket{tokens: float64(rpm), maxTokens: float64(rpm), refillRate: rate, last: time.Now()}
}

// Allow returns true if a request can proceed.
func (b *Bucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(b.last).Seconds()
	b.last = now
	b.tokens += elapsed * b.refillRate
	if b.tokens > b.maxTokens {
		b.tokens = b.maxTokens
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}
