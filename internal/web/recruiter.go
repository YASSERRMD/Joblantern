package web

import (
	"crypto/ed25519"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/yasserrmd/joblantern/internal/badge"
)

// BadgeIssuer wires the recruiter-side endpoints:
//
//	POST /api/v1/recruiter/badges  → issue a new badge from a verification
//	GET  /badge/{token}            → public verifier; returns the Claims
//	GET  /badge/{token}/svg        → small SVG you can <img src=...> embed
//
// Backed by an in-memory store for v1; a future PR persists to the
// recruiter_badges table (migration 0014).
type BadgeIssuer struct {
	IssuerURL string
	PrivKey   ed25519.PrivateKey
	PubKey    ed25519.PublicKey
	TTL       time.Duration

	mu      sync.Mutex
	revoked map[string]struct{}
}

// NewBadgeIssuer wires the routes onto r.
func NewBadgeIssuer(r chi.Router, issuerURL string, priv ed25519.PrivateKey, pub ed25519.PublicKey, ttl time.Duration) *BadgeIssuer {
	b := &BadgeIssuer{
		IssuerURL: issuerURL,
		PrivKey:   priv,
		PubKey:    pub,
		TTL:       ttl,
		revoked:   map[string]struct{}{},
	}
	r.Post("/api/v1/recruiter/badges", b.issue)
	r.Get("/badge/{token}", b.verify)
	r.Get("/badge/{token}/svg", b.svg)
	return b
}

type issueReq struct {
	OrgID          string `json:"org_id"`
	OrgName        string `json:"org_name"`
	VerificationID string `json:"verification_id"`
	Risk           string `json:"risk"`
	TrustLevel     string `json:"trust_level"`
}

type issueResp struct {
	Token  string       `json:"token"`
	Claims badge.Claims `json:"claims"`
}

func (b *BadgeIssuer) issue(w http.ResponseWriter, r *http.Request) {
	var req issueReq
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4<<10)).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	switch req.Risk {
	case "green", "yellow", "red":
	default:
		http.Error(w, "invalid risk", http.StatusBadRequest)
		return
	}
	if req.OrgID == "" || req.VerificationID == "" {
		http.Error(w, "org_id and verification_id required", http.StatusBadRequest)
		return
	}
	trust := req.TrustLevel
	if trust == "" {
		trust = "observed"
	}
	now := time.Now().UTC()
	c := badge.Claims{
		BadgeID:        uuid.NewString(),
		OrgID:          req.OrgID,
		OrgName:        req.OrgName,
		VerificationID: req.VerificationID,
		Risk:           req.Risk,
		IssuedAt:       now,
		ExpiresAt:      now.Add(b.TTL),
		Issuer:         b.IssuerURL,
		TrustLevel:     trust,
	}
	tok, err := badge.Encode(b.PrivKey, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, issueResp{Token: tok, Claims: c})
}

func (b *BadgeIssuer) verify(w http.ResponseWriter, r *http.Request) {
	tok := chi.URLParam(r, "token")
	c, err := badge.Decode(b.PubKey, tok)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"valid": "false", "error": err.Error(),
		})
		return
	}
	b.mu.Lock()
	_, revoked := b.revoked[c.BadgeID]
	b.mu.Unlock()
	if revoked {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"valid": false, "error": "revoked", "claims": c,
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"valid": true, "claims": c,
	})
}

// Revoke marks a badge id as revoked. Wired to an admin endpoint in
// a follow-up PR; exposed here so tests can exercise the path.
func (b *BadgeIssuer) Revoke(badgeID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.revoked[badgeID] = struct{}{}
}

// svg returns a minimal embeddable SVG showing the verdict colour.
func (b *BadgeIssuer) svg(w http.ResponseWriter, r *http.Request) {
	tok := chi.URLParam(r, "token")
	c, err := badge.Decode(b.PubKey, tok)
	if err != nil {
		http.Error(w, "invalid", http.StatusUnauthorized)
		return
	}
	colour := "#3da76a"
	switch c.Risk {
	case "yellow":
		colour = "#d6a23a"
	case "red":
		colour = "#d6452a"
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=600")
	_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="160" height="32" viewBox="0 0 160 32">` +
		`<rect width="160" height="32" rx="6" fill="` + colour + `"/>` +
		`<text x="80" y="20" font-family="-apple-system,system-ui,sans-serif" font-size="13" fill="#fff" text-anchor="middle" font-weight="600">Joblantern · ` + c.Risk + `</text>` +
		`</svg>`))
}
