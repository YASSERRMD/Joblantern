// Package federation implements Joblantern's peer-to-peer scam-signal
// exchange.
//
// Design summary:
//
//   - Each instance publishes a signed Manifest at
//     /.well-known/joblantern.json describing itself and its public key.
//   - Peers register each other (in the federated_peers table) by URL +
//     pubkey + trust_level.
//   - Outgoing signals are anonymised scam fingerprints: company name,
//     country, high-level pattern tags. **No user PII** ever leaves
//     the originating instance.
//   - Signals are signed with the originator's ed25519 private key.
//   - The receiving instance verifies the signature against the
//     stored peer pubkey and stores the signal with attribution.
//
// This package owns the wire types + the signing helpers; the actual
// HTTP surface lives in internal/web/federation.go.
package federation

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"
)

// Manifest is what an instance publishes at /.well-known/joblantern.json.
type Manifest struct {
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Version   string    `json:"version"`
	PubKeyHex string    `json:"pubkey_hex"`
	IssuedAt  time.Time `json:"issued_at"`
}

// Signal is one anonymised scam fingerprint forwarded between peers.
// No PII: company name and country only, plus the set of pattern codes
// the originator's risk engine assigned.
type Signal struct {
	OriginURL    string    `json:"origin_url"`
	CompanyName  string    `json:"company_name,omitempty"`
	Country      string    `json:"country,omitempty"`
	PatternCodes []string  `json:"pattern_codes,omitempty"`
	Verdict      string    `json:"verdict"` // "red" | "yellow"
	IssuedAt     time.Time `json:"issued_at"`
}

// Envelope is the on-the-wire representation: a Signal plus a
// detached ed25519 signature over its canonical JSON.
type Envelope struct {
	Signal       Signal `json:"signal"`
	SignatureB64 string `json:"signature_b64"`
}

// CanonicalBytes returns the deterministic JSON encoding signed for s.
// Map keys are alphabetical and string-array values are sorted so the
// same Signal always produces the same bytes regardless of source-side
// quirks.
func (s Signal) CanonicalBytes() ([]byte, error) {
	codes := append([]string(nil), s.PatternCodes...)
	sort.Strings(codes)
	m := map[string]any{
		"company_name":  s.CompanyName,
		"country":       s.Country,
		"issued_at":     s.IssuedAt.UTC().Format(time.RFC3339),
		"origin_url":    s.OriginURL,
		"pattern_codes": codes,
		"verdict":       s.Verdict,
	}
	return canonicalJSON(m)
}

// Sign produces an Envelope ready to POST to a peer.
func Sign(priv ed25519.PrivateKey, s Signal) (*Envelope, error) {
	if len(priv) == 0 {
		return nil, errors.New("federation: empty private key")
	}
	data, err := s.CanonicalBytes()
	if err != nil {
		return nil, err
	}
	sig := ed25519.Sign(priv, data)
	return &Envelope{Signal: s, SignatureB64: base64.StdEncoding.EncodeToString(sig)}, nil
}

// Verify confirms the envelope was signed by the holder of pub.
func Verify(pub ed25519.PublicKey, env *Envelope) error {
	if env == nil {
		return errors.New("federation: nil envelope")
	}
	data, err := env.Signal.CanonicalBytes()
	if err != nil {
		return err
	}
	sig, err := base64.StdEncoding.DecodeString(env.SignatureB64)
	if err != nil {
		return fmt.Errorf("federation: decode signature: %w", err)
	}
	if !ed25519.Verify(pub, data, sig) {
		return errors.New("federation: signature mismatch")
	}
	return nil
}

// Fingerprint is a content-hash of the signal (company + country +
// codes) used to dedupe identical signals from different peers.
func (s Signal) Fingerprint() string {
	h := sha256.New()
	h.Write([]byte(s.CompanyName))
	h.Write([]byte{0})
	h.Write([]byte(s.Country))
	codes := append([]string(nil), s.PatternCodes...)
	sort.Strings(codes)
	for _, c := range codes {
		h.Write([]byte{0})
		h.Write([]byte(c))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// canonicalJSON serialises with sorted keys recursively.
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
