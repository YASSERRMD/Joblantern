// Package snapshot exports annual archive snapshots to long-term
// storage. Snapshots are immutable, content-addressed, and
// integrity-hashed.
package snapshot

import "time"

// Spec drives a snapshot export.
type Spec struct {
	Year        int
	Region      string
	DestBucket  string
	StoreClass  string // "cold", "deep-archive", "tape"
}

// Result is the produced artifact.
type Result struct {
	URI         string
	SHA256      string
	Bytes       int64
	Snapshotted time.Time
}
