// Package badge issues and verifies signed "Verified by Joblantern"
// badges that recruiters and job boards can embed.
//
// Format (token): base64url( payload | "." | ed25519_signature(payload) )
// Payload is canonical JSON of Claims.
//
// The token is opaque to the embedder; the public verifier endpoint
// (/badge/<token>) is the canonical way to check it. The embedder
// renders a small HTML snippet pointing at that endpoint.
package badge

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"
)

// Claims is the verifiable content of a badge.
type Claims struct {
	BadgeID        string    `json:"badge_id"`
	OrgID          string    `json:"org_id"`
	OrgName        string    `json:"org_name"`
	VerificationID string    `json:"verification_id"`
	Risk           string    `json:"risk"` // green | yellow | red
	IssuedAt       time.Time `json:"issued_at"`
	ExpiresAt      time.Time `json:"expires_at"`
	Issuer         string    `json:"issuer"`      // joblantern instance URL
	TrustLevel     string    `json:"trust_level"` // vouched | observed
}

// Encode returns the signed token string.
func Encode(priv ed25519.PrivateKey, c Claims) (string, error) {
	if len(priv) == 0 {
		return "", errors.New("badge: empty private key")
	}
	payload, err := canonicalJSON(claimsMap(c))
	if err != nil {
		return "", err
	}
	sig := ed25519.Sign(priv, payload)
	tok := base64.RawURLEncoding.EncodeToString(payload) +
		"." + base64.RawURLEncoding.EncodeToString(sig)
	return tok, nil
}

// Decode verifies the token and returns its Claims.
func Decode(pub ed25519.PublicKey, tok string) (*Claims, error) {
	parts := splitDot(tok)
	if len(parts) != 2 {
		return nil, errors.New("badge: invalid token format")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("badge: decode payload: %w", err)
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("badge: decode signature: %w", err)
	}
	if !ed25519.Verify(pub, payload, sig) {
		return nil, errors.New("badge: signature mismatch")
	}
	var c Claims
	if err := json.Unmarshal(payload, &c); err != nil {
		return nil, fmt.Errorf("badge: decode claims: %w", err)
	}
	if time.Now().After(c.ExpiresAt) {
		return &c, errors.New("badge: expired")
	}
	return &c, nil
}

// claimsMap normalises a Claims into a sortable map for canonical JSON.
func claimsMap(c Claims) map[string]any {
	return map[string]any{
		"badge_id":        c.BadgeID,
		"org_id":          c.OrgID,
		"org_name":        c.OrgName,
		"verification_id": c.VerificationID,
		"risk":            c.Risk,
		"issued_at":       c.IssuedAt.UTC().Format(time.RFC3339),
		"expires_at":      c.ExpiresAt.UTC().Format(time.RFC3339),
		"issuer":          c.Issuer,
		"trust_level":     c.TrustLevel,
	}
}

func splitDot(s string) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}

func canonicalJSON(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return sortedMarshal(m)
}

func sortedMarshal(v any) ([]byte, error) {
	switch t := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		buf := []byte{'{'}
		for i, k := range keys {
			if i > 0 {
				buf = append(buf, ',')
			}
			kj, _ := json.Marshal(k)
			buf = append(buf, kj...)
			buf = append(buf, ':')
			vj, err := sortedMarshal(t[k])
			if err != nil {
				return nil, err
			}
			buf = append(buf, vj...)
		}
		buf = append(buf, '}')
		return buf, nil
	case []any:
		buf := []byte{'['}
		for i, el := range t {
			if i > 0 {
				buf = append(buf, ',')
			}
			eb, err := sortedMarshal(el)
			if err != nil {
				return nil, err
			}
			buf = append(buf, eb...)
		}
		buf = append(buf, ']')
		return buf, nil
	default:
		return json.Marshal(v)
	}
}
