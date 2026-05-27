// Package audit appends an immutable trail of regulator-side actions.
// Entries are hash-chained so any tampering is detectable.
package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"
)

// Entry is one audit row.
type Entry struct {
	At        time.Time `json:"at"`
	Regulator string    `json:"regulator"`
	Subject   string    `json:"subject"`
	Action    string    `json:"action"`
	Detail    string    `json:"detail,omitempty"`
	PrevHash  string    `json:"prev_hash"`
	Hash      string    `json:"hash"`
}

// Log is an in-memory hash-chained audit log. Production replaces
// this with a database-backed implementation that flushes a Merkle
// root to public timestamping (e.g., the transparency log from
// Phase 9).
type Log struct {
	mu      sync.Mutex
	entries []Entry
	last    string
}

// Append records an action and returns the new tip hash.
func (l *Log) Append(e Entry) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e.At = time.Now().UTC()
	e.PrevHash = l.last
	body, err := json.Marshal(struct {
		At        time.Time `json:"at"`
		Regulator string    `json:"regulator"`
		Subject   string    `json:"subject"`
		Action    string    `json:"action"`
		Detail    string    `json:"detail,omitempty"`
		PrevHash  string    `json:"prev_hash"`
	}{e.At, e.Regulator, e.Subject, e.Action, e.Detail, e.PrevHash})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(body)
	e.Hash = hex.EncodeToString(sum[:])
	l.entries = append(l.entries, e)
	l.last = e.Hash
	return e.Hash, nil
}

// Entries returns a snapshot of the log.
func (l *Log) Entries() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}
