// Package sandbox provisions a sandboxed Postgres replica for a
// research partnership. The sandbox is a logical-replication target
// of the production primary, scrubbed of submitter identities and
// raw contact fields before any researcher account connects.
package sandbox

import (
	"errors"
	"time"
)

// Spec is what the partnership requested.
type Spec struct {
	PartnerID   string
	Region      string // residency-compliant
	StartAt     time.Time
	EndAt       time.Time
	Tier        string // "B" or "C"
	Fields      []string
}

// Provisioning is the lifecycle status.
type Provisioning struct {
	Spec       Spec
	State      string // "pending", "ready", "destroyed"
	StateAt    time.Time
	Hostname   string
	Database   string
}

// Validate enforces the minimum invariants.
func (s Spec) Validate() error {
	if s.PartnerID == "" {
		return errors.New("partner required")
	}
	if s.StartAt.IsZero() || s.EndAt.IsZero() {
		return errors.New("window required")
	}
	if !s.EndAt.After(s.StartAt) {
		return errors.New("end must be after start")
	}
	if s.Tier != "B" && s.Tier != "C" {
		return errors.New("tier must be B or C")
	}
	return nil
}
