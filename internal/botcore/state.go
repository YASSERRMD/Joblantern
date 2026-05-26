package botcore

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// State is a per-chat conversation state. v1 is intentionally tiny:
// the bot either has a pending submission half-built or it does not.
type State struct {
	ChatID     string
	Submission Submission
	UpdatedAt  time.Time
	LastID     string // most recent verification id we kicked off
}

// Sessions is an in-memory chat-keyed state store.
type Sessions struct {
	mu sync.Mutex
	m  map[string]*State
}

// NewSessions returns an empty store.
func NewSessions() *Sessions {
	return &Sessions{m: map[string]*State{}}
}

// Get returns the state for chat (creating an empty one if absent).
func (s *Sessions) Get(chat string) *State {
	s.mu.Lock()
	defer s.mu.Unlock()
	st := s.m[chat]
	if st == nil {
		st = &State{ChatID: chat, UpdatedAt: time.Now()}
		s.m[chat] = st
	}
	return st
}

// Reset clears chat state (used by /forget).
func (s *Sessions) Reset(chat string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, chat)
}

// IPRateLimiter — adapter-agnostic per-key token bucket.
type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]int
	resets  map[string]time.Time
	limit   int
	window  time.Duration
}

// NewRateLimiter allows `limit` events per `window` per key.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		buckets: map[string]int{},
		resets:  map[string]time.Time{},
		limit:   limit,
		window:  window,
	}
}

// Allow returns true if the key has tokens remaining.
func (l *RateLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if r, ok := l.resets[key]; !ok || now.After(r) {
		l.resets[key] = now.Add(l.window)
		l.buckets[key] = l.limit
	}
	if l.buckets[key] <= 0 {
		return false
	}
	l.buckets[key]--
	return true
}

// FormatVerdict turns a Record into a Telegram/WhatsApp-friendly reply.
// Kept transport-neutral so every adapter renders the same content.
func FormatVerdict(rec *Record, viewURL string) string {
	if rec == nil {
		return "no record"
	}
	if rec.Verdict == nil {
		return fmt.Sprintf("Status: %s. Still working — try /status %s in a moment.", rec.Status, rec.ID)
	}
	v := rec.Verdict
	risk := strings.ToUpper(v.OverallRisk)
	var b strings.Builder
	switch v.OverallRisk {
	case "red":
		b.WriteString("🚨 ")
	case "yellow":
		b.WriteString("⚠️ ")
	default:
		b.WriteString("✅ ")
	}
	fmt.Fprintf(&b, "%s · confidence %.0f%%\n", risk, v.Confidence*100)
	if len(v.Reasons) > 0 {
		b.WriteString("\nReasons:\n")
		for i, r := range v.Reasons {
			if i >= 5 {
				break
			}
			fmt.Fprintf(&b, "  • %s\n", r)
		}
	}
	if viewURL != "" {
		fmt.Fprintf(&b, "\nFull report: %s", viewURL)
	}
	return b.String()
}
