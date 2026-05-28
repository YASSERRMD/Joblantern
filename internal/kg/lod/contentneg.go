// Package lod implements Linked Open Data content negotiation. A GET
// on an entity IRI returns HTML or RDF depending on the Accept
// header.
package lod

import (
	"net/http"
	"strings"
)

// Format is the selected response format.
type Format string

const (
	FormatHTML   Format = "text/html"
	FormatTurtle Format = "text/turtle"
	FormatJSONLD Format = "application/ld+json"
)

// Negotiate picks a format from an Accept header.
func Negotiate(accept string) Format {
	a := strings.ToLower(accept)
	switch {
	case strings.Contains(a, "application/ld+json"):
		return FormatJSONLD
	case strings.Contains(a, "text/turtle"):
		return FormatTurtle
	default:
		return FormatHTML
	}
}

// Handler is the LOD content-negotiating handler.
type Handler struct {
	HTML   func(w http.ResponseWriter, r *http.Request)
	Turtle func(w http.ResponseWriter, r *http.Request)
	JSONLD func(w http.ResponseWriter, r *http.Request)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch Negotiate(r.Header.Get("Accept")) {
	case FormatJSONLD:
		h.JSONLD(w, r)
	case FormatTurtle:
		h.Turtle(w, r)
	default:
		h.HTML(w, r)
	}
}
