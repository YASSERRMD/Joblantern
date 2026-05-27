// Package residency selects the region a tenant's data must live in.
// We support EU, US, and APAC and never cross a region for tenant
// data even if it would be operationally easier.
package residency

import "errors"

// Region is the canonical id.
type Region string

const (
	EU   Region = "eu"
	US   Region = "us"
	APAC Region = "apac"
)

// EndpointForRegion returns the public ingress endpoint for a region.
// The caller validates that the active request hit the right region.
func EndpointForRegion(r Region) (string, error) {
	switch r {
	case EU:
		return "https://eu.joblantern.org", nil
	case US:
		return "https://us.joblantern.org", nil
	case APAC:
		return "https://ap.joblantern.org", nil
	}
	return "", errors.New("unknown region")
}

// ValidateCountry returns the default region for a country code.
// Tenants may override but must justify the override in onboarding.
func ValidateCountry(iso2 string) Region {
	switch iso2 {
	case "DE", "FR", "ES", "IT", "NL", "BE", "IE", "PT", "SE", "FI", "DK", "AT":
		return EU
	case "US", "CA":
		return US
	}
	return APAC
}
