// Package openapi generates an OpenAPI 3.1 specification from the chi
// router used by the public researcher API. Routes register their
// metadata (summary, parameters, responses) through Register, and the
// final spec is rendered via Render.
package openapi

import (
	"encoding/json"
	"sort"
	"sync"
)

// Operation describes a single REST operation in OpenAPI 3.1 form.
type Operation struct {
	Method      string              `json:"-"`
	Path        string              `json:"-"`
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	Tags        []string            `json:"tags,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty"`
}

// Parameter is an OpenAPI parameter object.
type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Required    bool   `json:"required,omitempty"`
	Description string `json:"description,omitempty"`
	Schema      Schema `json:"schema"`
}

// Schema is a minimal OpenAPI schema reference.
type Schema struct {
	Type   string `json:"type,omitempty"`
	Format string `json:"format,omitempty"`
	Ref    string `json:"$ref,omitempty"`
}

// Response is an OpenAPI response object.
type Response struct {
	Description string `json:"description"`
}

// Registry stores routes for spec generation.
type Registry struct {
	mu  sync.Mutex
	ops []Operation
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry { return &Registry{} }

// Register adds an operation.
func (r *Registry) Register(op Operation) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ops = append(r.ops, op)
}

// Render returns an OpenAPI 3.1 document as JSON bytes.
func (r *Registry) Render(title, version string) ([]byte, error) {
	r.mu.Lock()
	ops := append([]Operation(nil), r.ops...)
	r.mu.Unlock()
	sort.Slice(ops, func(i, j int) bool {
		if ops[i].Path == ops[j].Path {
			return ops[i].Method < ops[j].Method
		}
		return ops[i].Path < ops[j].Path
	})
	paths := map[string]map[string]Operation{}
	for _, op := range ops {
		bucket, ok := paths[op.Path]
		if !ok {
			bucket = map[string]Operation{}
			paths[op.Path] = bucket
		}
		bucket[op.Method] = op
	}
	doc := map[string]any{
		"openapi": "3.1.0",
		"info": map[string]any{
			"title":   title,
			"version": version,
		},
		"paths": paths,
	}
	return json.MarshalIndent(doc, "", "  ")
}
