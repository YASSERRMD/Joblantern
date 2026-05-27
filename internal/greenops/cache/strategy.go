// Package cache documents the caching strategy designed to minimise
// redundant MCP calls. Every cache hit is a saved network call, which
// is a saved request to a remote service running on a different grid.
package cache

import "time"

// Policy enumerates the per-call-class TTLs we ship.
type Policy struct {
	Class string
	TTL   time.Duration
	Notes string
}

// Defaults captures the policy in one place so the green-ops report
// can compare aggregate savings against the LLM cost.
func Defaults() []Policy {
	return []Policy{
		{Class: "registry-lookup", TTL: 7 * 24 * time.Hour, Notes: "Company registry changes infrequently."},
		{Class: "domain-whois", TTL: 24 * time.Hour, Notes: "WHOIS rarely changes within a day."},
		{Class: "salary-bands", TTL: 30 * 24 * time.Hour, Notes: "Statistics update at most monthly."},
		{Class: "accreditation", TTL: 30 * 24 * time.Hour, Notes: "Accreditation rolls update infrequently."},
		{Class: "rental-listings", TTL: 1 * time.Hour, Notes: "Listings expire quickly."},
		{Class: "scam-db-hash", TTL: 6 * time.Hour, Notes: "Balance freshness against load."},
	}
}
