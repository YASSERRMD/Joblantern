// Package queue abstracts the verification job queue. The first
// supported backend is River (Apache 2.0) over Postgres; Asynq over
// Redis is an alternative for cache-heavy regions.
package queue

import (
	"context"
	"errors"
	"time"
)

// Job is the minimal job envelope.
type Job struct {
	ID        string
	Kind      string
	Payload   []byte
	Attempts  int
	EnqueueAt time.Time
}

// Queue is the abstraction.
type Queue interface {
	Enqueue(ctx context.Context, j Job) error
	Dequeue(ctx context.Context) (*Job, error)
	Complete(ctx context.Context, jobID string) error
	Retry(ctx context.Context, jobID string, backoff time.Duration) error
}

// MaxAttempts is the global retry cap.
const MaxAttempts = 6

// NextBackoff returns the suggested backoff for an attempt.
func NextBackoff(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	if attempt > 8 {
		attempt = 8
	}
	return time.Duration(1<<attempt) * time.Second
}

// ErrEmpty is returned when there is no job to dequeue right now.
var ErrEmpty = errors.New("queue empty")
