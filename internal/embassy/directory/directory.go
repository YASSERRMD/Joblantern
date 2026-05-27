// Package directory holds national embassy / hotline directories.
// Entries are loaded from a YAML file at kiosk boot so deployments
// can update contacts without rebuilding the binary.
package directory

import (
	"errors"
	"strings"
)

// Contact is one entry in the embassy directory.
type Contact struct {
	Country     string
	Agency      string
	Hotline     string
	OfficeEmail string
	OfficeURL   string
	Languages   []string
}

// Directory is the in-memory lookup table.
type Directory struct {
	byCountry map[string][]Contact
}

// New constructs an empty directory.
func New() *Directory { return &Directory{byCountry: map[string][]Contact{}} }

// Add inserts a contact.
func (d *Directory) Add(c Contact) error {
	if strings.TrimSpace(c.Country) == "" || strings.TrimSpace(c.Hotline) == "" {
		return errors.New("country and hotline required")
	}
	key := strings.ToUpper(c.Country)
	d.byCountry[key] = append(d.byCountry[key], c)
	return nil
}

// Lookup returns all contacts for a country (ISO-2).
func (d *Directory) Lookup(country string) []Contact {
	return d.byCountry[strings.ToUpper(country)]
}
