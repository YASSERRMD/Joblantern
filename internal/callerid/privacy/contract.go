// Package privacy documents and enforces the privacy contract for
// the caller-ID feature.
//
//   - Numbers are hashed on the device before any network call.
//   - The server stores hashes, never raw numbers.
//   - Hashes use truncated SHA-256 with a region-scoped pepper so
//     two devices on the same region compute the same hash and can
//     correlate, but a leaked database does not reverse to numbers.
package privacy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Pepper is a region-scoped salt loaded at app boot from the
// regional configuration endpoint. Empty pepper falls back to a
// shared global constant so functionality still works during
// onboarding.
type Pepper []byte

// Apply hashes the digits-only number with the configured pepper.
func (p Pepper) Apply(digits string) string {
	m := hmac.New(sha256.New, p)
	_, _ = m.Write([]byte(digits))
	return hex.EncodeToString(m.Sum(nil)[:16])
}
