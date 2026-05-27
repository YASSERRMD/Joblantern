// Package triplestore is the abstraction over the RDF triple store
// backing Joblantern's knowledge graph. The default is a small
// Postgres-backed store; large deployments swap in Apache Jena
// Fuseki as a sidecar service.
package triplestore

import (
	"context"
	"errors"
)

// Triple is one RDF triple.
type Triple struct {
	Subject   string // IRI
	Predicate string // IRI
	Object    string // IRI or literal
	Graph     string // named graph, optional
}

// Store is the abstraction.
type Store interface {
	Put(ctx context.Context, t Triple) error
	PutBatch(ctx context.Context, ts []Triple) error
	Delete(ctx context.Context, t Triple) error
	Query(ctx context.Context, sparql string) ([]map[string]string, error)
}

// ErrUnsupported is returned by stores that decline a feature (e.g.
// the embedded Postgres store does not support full SPARQL update).
var ErrUnsupported = errors.New("operation not supported by this store")
