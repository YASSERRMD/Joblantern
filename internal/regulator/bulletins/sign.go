// Package bulletins issues publicly verifiable signed bulletins from
// regulators. Each bulletin is served at /bulletins/<regulator>/<id>
// and includes a detached ed25519 signature plus the regulator's
// public key fingerprint.
package bulletins

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

// Bulletin is one signed public statement.
type Bulletin struct {
	ID          string    `json:"id"`
	Regulator   string    `json:"regulator"`
	IssuedAt    time.Time `json:"issued_at"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Signature   string    `json:"signature"`
	KeyFingerprint string `json:"key_fingerprint"`
}

// canonical returns the canonical JSON used as the signing payload.
// The Signature field is excluded from the payload.
func (b *Bulletin) canonical() ([]byte, error) {
	copy := *b
	copy.Signature = ""
	return json.Marshal(copy)
}

// Sign produces a detached signature using the regulator private key.
func (b *Bulletin) Sign(priv ed25519.PrivateKey) error {
	payload, err := b.canonical()
	if err != nil {
		return err
	}
	sig := ed25519.Sign(priv, payload)
	b.Signature = hex.EncodeToString(sig)
	pub := priv.Public().(ed25519.PublicKey)
	fp := sha256.Sum256(pub)
	b.KeyFingerprint = hex.EncodeToString(fp[:])
	return nil
}

// Verify checks the bulletin against the supplied public key.
func (b Bulletin) Verify(pub ed25519.PublicKey) error {
	payload, err := b.canonical()
	if err != nil {
		return err
	}
	sig, err := hex.DecodeString(b.Signature)
	if err != nil {
		return errors.New("malformed signature")
	}
	if !ed25519.Verify(pub, payload, sig) {
		return errors.New("signature does not verify")
	}
	return nil
}
