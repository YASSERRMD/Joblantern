// Package entity resolves entities (phones, emails, addresses,
// director names) into canonical node identifiers used by the graph.
//
// Resolution is deterministic so two ingestion paths see the same
// node id for the same underlying value.
package entity

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// Kind identifies the entity type.
type Kind string

const (
	KindPhone    Kind = "phone"
	KindEmail    Kind = "email"
	KindAddress  Kind = "address"
	KindDirector Kind = "director"
	KindPayment  Kind = "payment"
)

// Node is the canonical representation.
type Node struct {
	ID    string
	Kind  Kind
	Label string
}

// Resolve normalises the raw value and produces a stable node id.
func Resolve(kind Kind, raw string) Node {
	n := normalise(kind, raw)
	sum := sha256.Sum256([]byte(string(kind) + ":" + n))
	return Node{ID: hex.EncodeToString(sum[:12]), Kind: kind, Label: n}
}

func normalise(kind Kind, raw string) string {
	raw = strings.TrimSpace(raw)
	switch kind {
	case KindPhone:
		var b strings.Builder
		for _, r := range raw {
			if r >= '0' && r <= '9' {
				b.WriteRune(r)
			}
		}
		return b.String()
	case KindEmail:
		return strings.ToLower(raw)
	case KindAddress, KindDirector, KindPayment:
		return strings.ToLower(strings.Join(strings.Fields(raw), " "))
	}
	return raw
}
