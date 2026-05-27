// Package apikey issues, verifies, and rate-limits per-tenant API
// keys. Keys are stored as their argon2id digest.
package apikey

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Key is a tenant-scoped credential.
type Key struct {
	ID         string
	TenantID   string
	Prefix     string // first 8 chars of the public key, for ops
	HashedKey  []byte
	Scopes     []string
	RatePerMin int
}

// Fingerprint returns a stable short identifier (not the secret).
func (k Key) Fingerprint() string {
	sum := sha256.Sum256(k.HashedKey)
	return hex.EncodeToString(sum[:8])
}

// Compare is a constant-time comparator for HMAC-derived material.
// In production secret comparison uses argon2id verify; this helper
// is for the prefix lookup path.
func Compare(a, b []byte) bool { return hmac.Equal(a, b) }
