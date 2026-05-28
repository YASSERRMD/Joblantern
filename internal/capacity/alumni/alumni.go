// Package alumni tracks operators who have moved on but stay
// connected to the network.
package alumni

import "time"

// Member is one alumni record.
type Member struct {
	ID            string
	OperatorID    string
	LeftAt        time.Time
	Reason        string
	StayConnected bool
}

// Active reports whether the alumnus is still subscribed to alumni
// communications.
func Active(m Member) bool { return m.StayConnected }
