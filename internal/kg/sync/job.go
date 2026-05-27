// Package sync replays Postgres rows into the triple store on a
// cadence appropriate to the deployment.
package sync

import "time"

// Cursor is the last-synced position.
type Cursor struct {
	LastVerifiedAt time.Time
	LastID         string
}

// BatchSize is the upper bound per sync iteration.
const BatchSize = 1000

// Interval is the default sync cadence.
const Interval = 5 * time.Minute

// NextCursor advances the cursor given a batch of rows. Callers
// supply the max(verified_at, id) of the consumed batch.
func NextCursor(cur Cursor, maxAt time.Time, maxID string) Cursor {
	if maxAt.After(cur.LastVerifiedAt) {
		cur.LastVerifiedAt = maxAt
		cur.LastID = maxID
	}
	return cur
}
