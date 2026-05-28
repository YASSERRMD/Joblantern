// Package operator extends the Phase 27 recruiter API for marketplace
// operators. Operators can subscribe to scam-listing notifications
// keyed on their own listing IDs.
package operator

import "time"

// Subscription is one operator's webhook subscription for incidents
// touching their listings.
type Subscription struct {
	OperatorID string
	URL        string
	Secret     []byte
	CreatedAt  time.Time
}

// Notification is the wire form delivered to operator endpoints.
type Notification struct {
	IncidentID string    `json:"incident_id"`
	ListingID  string    `json:"listing_id"`
	Band       string    `json:"band"`
	Confidence float64   `json:"confidence"`
	At         time.Time `json:"at"`
}
