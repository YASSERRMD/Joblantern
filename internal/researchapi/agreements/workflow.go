// Package agreements implements an automated, DocuSign-style data-use
// agreement workflow without depending on any specific e-signature
// vendor. The flow is: request -> review -> applicant signs (ed25519)
// -> Joblantern counter-signs -> token issued.
package agreements

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
)

// Status is the lifecycle of an agreement.
type Status string

const (
	StatusDraft         Status = "draft"
	StatusAwaitingSign  Status = "awaiting_signature"
	StatusCounterSigned Status = "counter_signed"
	StatusRevoked       Status = "revoked"
)

// Agreement is a researcher data-use agreement.
type Agreement struct {
	ID                string
	ResearcherEmail   string
	Organization      string
	Tier              string
	CreatedAt         time.Time
	Status            Status
	BodySHA256        string
	ResearcherSig     []byte
	JoblanternSig     []byte
	ResearcherPubKey  ed25519.PublicKey
	JoblanternPubKey  ed25519.PublicKey
	Body              string
}

// Sign hashes the body, validates the researcher signature and moves
// the agreement to awaiting counter-signature.
func (a *Agreement) Sign(sig []byte) error {
	sum := sha256.Sum256([]byte(a.Body))
	a.BodySHA256 = hex.EncodeToString(sum[:])
	if !ed25519.Verify(a.ResearcherPubKey, sum[:], sig) {
		return errors.New("researcher signature does not verify")
	}
	a.ResearcherSig = sig
	a.Status = StatusAwaitingSign
	return nil
}

// CounterSign produces the joblantern signature using the supplied
// private key and finalizes the agreement.
func (a *Agreement) CounterSign(priv ed25519.PrivateKey) {
	sum := sha256.Sum256([]byte(a.Body))
	a.JoblanternSig = ed25519.Sign(priv, sum[:])
	a.JoblanternPubKey = priv.Public().(ed25519.PublicKey)
	a.Status = StatusCounterSigned
}
