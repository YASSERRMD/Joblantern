package web

import (
	"crypto/ed25519"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/yasserrmd/joblantern/internal/federation"
)

// PeerLookup resolves a peer's public key by its base URL.
type PeerLookup func(url string) (ed25519.PublicKey, bool)

// FederationHandler exposes:
//   - GET  /.well-known/joblantern.json  → signed Manifest
//   - POST /a2a/signal                   → ingest a signed Envelope
type FederationHandler struct {
	Name      string
	URL       string
	Version   string
	PubKey    ed25519.PublicKey
	Lookup    PeerLookup
	mu        sync.Mutex
	signals   []federation.Signal // in-memory ring for v1; future PR persists
	maxStored int
}

// NewFederation wires the federation routes onto r.
func NewFederation(r chi.Router, name, url, version string, pub ed25519.PublicKey, lookup PeerLookup) *FederationHandler {
	h := &FederationHandler{
		Name: name, URL: url, Version: version, PubKey: pub, Lookup: lookup,
		maxStored: 1024,
	}
	r.Get("/.well-known/joblantern.json", h.manifest)
	r.Post("/a2a/signal", h.ingest)
	r.Get("/a2a/recent", h.recent)
	return h
}

func (h *FederationHandler) manifest(w http.ResponseWriter, _ *http.Request) {
	m := federation.Manifest{
		Name:      h.Name,
		URL:       h.URL,
		Version:   h.Version,
		PubKeyHex: hex(h.PubKey),
		IssuedAt:  time.Now().UTC(),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(m)
}

func (h *FederationHandler) ingest(w http.ResponseWriter, r *http.Request) {
	var env federation.Envelope
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10)).Decode(&env); err != nil {
		http.Error(w, "bad envelope", http.StatusBadRequest)
		return
	}
	pub, ok := h.Lookup(env.Signal.OriginURL)
	if !ok {
		http.Error(w, "unknown peer", http.StatusForbidden)
		return
	}
	if err := federation.Verify(pub, &env); err != nil {
		http.Error(w, "signature: "+err.Error(), http.StatusUnauthorized)
		return
	}
	h.store(env.Signal)
	w.WriteHeader(http.StatusAccepted)
}

func (h *FederationHandler) recent(w http.ResponseWriter, _ *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]federation.Signal, len(h.signals))
	copy(out, h.signals)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func (h *FederationHandler) store(s federation.Signal) {
	h.mu.Lock()
	defer h.mu.Unlock()
	// Dedupe by fingerprint within the ring buffer.
	fp := s.Fingerprint()
	for _, existing := range h.signals {
		if existing.Fingerprint() == fp {
			return
		}
	}
	h.signals = append(h.signals, s)
	if len(h.signals) > h.maxStored {
		h.signals = h.signals[len(h.signals)-h.maxStored:]
	}
}

const hexChars = "0123456789abcdef"

func hex(b []byte) string {
	out := make([]byte, len(b)*2)
	for i, c := range b {
		out[i*2] = hexChars[c>>4]
		out[i*2+1] = hexChars[c&0x0f]
	}
	return string(out)
}
