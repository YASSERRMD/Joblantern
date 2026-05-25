// Package registry defines the abstract business-registry provider
// interface and the data types Joblantern exposes through the
// mcp-registry server. Concrete providers (OpenCorporates, UK
// Companies House, etc.) live under internal/registry/<name>/.
package registry

import (
	"context"
	"errors"
	"time"
)

// Common errors returned by Provider implementations.
var (
	ErrNotFound       = errors.New("registry: company not found")
	ErrRateLimited    = errors.New("registry: rate limited")
	ErrTokenInvalid   = errors.New("registry: token invalid")
	ErrJurisdiction   = errors.New("registry: jurisdiction unknown")
	ErrNotImplemented = errors.New("registry: not implemented")
)

// Match is a search hit.
type Match struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Jurisdiction      string    `json:"jurisdiction,omitempty"`
	Status            string    `json:"status,omitempty"`
	IncorporationDate time.Time `json:"incorporation_date,omitempty"`
	RegistryURL       string    `json:"registry_url,omitempty"`
}

// Company is the full record for a single company.
type Company struct {
	Match
	RegisteredAddress string    `json:"registered_address,omitempty"`
	Officers          []Officer `json:"officers,omitempty"`
}

// Officer is a director / officer record.
type Officer struct {
	Name      string    `json:"name"`
	Position  string    `json:"position,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
}

// Provider is the interface every registry backend implements.
type Provider interface {
	// Name returns a short stable identifier (e.g. "opencorporates").
	Name() string

	// LookupByName searches by company name, optionally scoped to a
	// jurisdiction (ISO-3166-style codes such as "ae", "gb", "us_de").
	LookupByName(ctx context.Context, name, jurisdiction string, limit int) ([]Match, error)

	// Get fetches the full record by the provider-specific id.
	Get(ctx context.Context, id string) (*Company, error)
}
