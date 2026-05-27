// Package throttle gates automatic submissions. By policy no takedown
// ships without an explicit human approval recorded in the case
// docket.
package throttle

import (
	"errors"
	"sync"
	"time"
)

// Decision is the human approval payload.
type Decision struct {
	ReviewerID string
	At         time.Time
	Notes      string
}

// Queue holds packets pending approval.
type Queue struct {
	mu       sync.Mutex
	pending  map[string]Decision // packetID -> empty until approved
	approved map[string]Decision
}

// New constructs an empty queue.
func New() *Queue { return &Queue{pending: map[string]Decision{}, approved: map[string]Decision{}} }

// Enqueue records a packet awaiting approval.
func (q *Queue) Enqueue(packetID string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.approved[packetID]; ok {
		return
	}
	if _, ok := q.pending[packetID]; ok {
		return
	}
	q.pending[packetID] = Decision{}
}

// Approve marks a packet ready to send.
func (q *Queue) Approve(packetID string, d Decision) error {
	if d.ReviewerID == "" {
		return errors.New("reviewer required")
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.pending[packetID]; !ok {
		return errors.New("not pending")
	}
	delete(q.pending, packetID)
	d.At = time.Now().UTC()
	q.approved[packetID] = d
	return nil
}

// Ready returns the approved packet ids.
func (q *Queue) Ready() []string {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]string, 0, len(q.approved))
	for id := range q.approved {
		out = append(out, id)
	}
	return out
}
