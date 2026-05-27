// Package feedback receives "this isn't a scam, I know this number"
// signals from the device. Repeated independent dismissals from
// trusted users de-weight a hash.
package feedback

import "time"

// Signal is the wire form of a false-positive feedback event.
type Signal struct {
	PhoneHash string    `json:"phone_hash"`
	UserToken string    `json:"user_token"` // pseudonymous, not the user id
	At        time.Time `json:"at"`
	Note      string    `json:"note,omitempty"`
}

// QuorumThreshold is the number of distinct trusted users whose
// dismissal moves a hash from red to yellow pending council review.
const QuorumThreshold = 5

// QuorumReached returns true if the supplied signals are from at
// least QuorumThreshold distinct trusted users.
func QuorumReached(s []Signal) bool {
	seen := map[string]struct{}{}
	for _, x := range s {
		if x.UserToken == "" {
			continue
		}
		seen[x.UserToken] = struct{}{}
	}
	return len(seen) >= QuorumThreshold
}
