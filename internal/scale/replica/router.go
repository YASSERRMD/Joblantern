// Package replica routes read-only evidence queries to a Postgres
// read replica. Write paths always go to the primary.
package replica

import "context"

// Pool is the minimum interface a DB pool needs to expose so the
// router can dispatch.
type Pool interface {
	Exec(ctx context.Context, sql string, args ...any) error
	Query(ctx context.Context, sql string, args ...any) ([]map[string]any, error)
}

// Router picks the right pool for the operation.
type Router struct {
	Primary  Pool
	Replicas []Pool
	next     int
}

// Read returns one of the replicas if available, else falls back to
// the primary.
func (r *Router) Read() Pool {
	if len(r.Replicas) == 0 {
		return r.Primary
	}
	p := r.Replicas[r.next%len(r.Replicas)]
	r.next++
	return p
}

// Write always returns the primary.
func (r *Router) Write() Pool { return r.Primary }
