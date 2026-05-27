package lookup

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// HashPhone normalises a phone number (digits only) and returns the
// truncated SHA-256 hex digest. The server stores hashes only — raw
// numbers never leave the device.
func HashPhone(raw string) string {
	var sb strings.Builder
	for _, r := range raw {
		if r >= '0' && r <= '9' {
			sb.WriteRune(r)
		}
	}
	sum := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(sum[:16])
}
