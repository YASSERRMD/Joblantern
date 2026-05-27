// Package regulator implements the integration layer with official
// labor ministries, embassies, and consumer-protection regulators.
//
// Regulator accounts are a high-trust tier and are verified through
// two channels:
//
//  1. The applicant must register from an official domain whose MX
//     and TXT records resolve to the published government registrar.
//  2. The applicant must upload a signed letter from a published
//     point of contact; the signature is verified offline by the
//     trust & safety council.
package regulator

import (
	"errors"
	"net"
	"strings"
	"time"
)

// Account is a verified regulator principal.
type Account struct {
	ID            string
	Country       string
	Agency        string
	OfficialEmail string
	OfficialDomain string
	LetterSHA256  string
	VerifiedAt    time.Time
	MTLSFingerprint string
}

// Resolver wraps the small surface of DNS we need.
type Resolver interface {
	LookupMX(domain string) ([]*net.MX, error)
	LookupTXT(domain string) ([]string, error)
}

// VerifyDomain confirms the supplied domain advertises the regulator
// registrar TXT record.
func VerifyDomain(r Resolver, domain, registrarMarker string) error {
	domain = strings.ToLower(strings.TrimSpace(domain))
	if domain == "" {
		return errors.New("empty domain")
	}
	if _, err := r.LookupMX(domain); err != nil {
		return errors.New("domain has no mail server")
	}
	txts, err := r.LookupTXT(domain)
	if err != nil {
		return errors.New("TXT lookup failed")
	}
	for _, t := range txts {
		if strings.Contains(t, registrarMarker) {
			return nil
		}
	}
	return errors.New("registrar TXT marker not found")
}
