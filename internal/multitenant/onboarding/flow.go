// Package onboarding handles tenant signup with a light-touch KYC
// check appropriate for an NGO-led service. The default flow lets a
// new NGO self-serve to first verdict in under an hour.
package onboarding

import (
	"errors"
	"strings"
)

// Application is the signup payload.
type Application struct {
	Slug          string
	DisplayName   string
	ContactEmail  string
	Country       string
	Mission       string
	NGORegistry   string
	NGORegistryID string
}

// KYC level returns the required documentation tier.
type KYC string

const (
	KYCLite KYC = "lite"
	KYCFull KYC = "full"
)

// Required returns the KYC level for the country. NGO-led signups in
// permissive jurisdictions only need a lite verification; commercial
// recruiters always need full KYC.
func Required(country, role string) KYC {
	if strings.EqualFold(role, "commercial-recruiter") {
		return KYCFull
	}
	return KYCLite
}

// Validate enforces the signup minimums.
func (a Application) Validate() error {
	if strings.TrimSpace(a.Slug) == "" {
		return errors.New("slug required")
	}
	if !strings.Contains(a.ContactEmail, "@") {
		return errors.New("contact email required")
	}
	if a.NGORegistry != "" && a.NGORegistryID == "" {
		return errors.New("registry id required when registry is set")
	}
	return nil
}
