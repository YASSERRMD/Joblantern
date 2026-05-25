// Package web wires the public HTTP surface — JSON API in Phase 13,
// templ-rendered UI later in Phase 15.
package web

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/yasserrmd/joblantern/internal/agent"
)

// APIHandler binds /api/v1 routes to a Router.
type APIHandler struct {
	Store        agent.Store
	Orchestrator *agent.Orchestrator
}

// NewAPIHandler wires routes onto r and returns h for testability.
func NewAPIHandler(r chi.Router, store agent.Store, orch *agent.Orchestrator) *APIHandler {
	h := &APIHandler{Store: store, Orchestrator: orch}
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/verify", h.postVerify)
		r.Get("/verifications", h.listVerifications)
		r.Get("/verifications/{id}", h.getVerification)
	})
	return h
}

func (h *APIHandler) postVerify(w http.ResponseWriter, r *http.Request) {
	var sub agent.Submission
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&sub); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json: " + err.Error()})
		return
	}
	sub.Jurisdiction = strings.ToUpper(strings.TrimSpace(sub.Jurisdiction))

	id, err := h.Store.Create(r.Context(), sub)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Run asynchronously so the client gets the id immediately.
	go h.RunAndStore(context.Background(), id, sub)

	writeJSON(w, http.StatusAccepted, map[string]string{"verification_id": id})
}

// RunAndStore exposes the async run for UI handlers that share the same APIHandler.
func (h *APIHandler) RunAndStore(ctx context.Context, id string, sub agent.Submission) {
	rec, _ := h.Store.Get(ctx, id)
	if rec != nil {
		rec.Status = "running"
		_ = h.Store.Save(ctx, rec)
	}
	sub.ID = id
	verdict := h.Orchestrator.Run(ctx, sub)
	if rec == nil {
		rec = &agent.Record{ID: id, Submission: sub}
	}
	rec.Status = "completed"
	rec.Verdict = &verdict
	_ = h.Store.Save(ctx, rec)
}

func (h *APIHandler) getVerification(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	rec, err := h.Store.Get(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if rec == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, rec)
}

func (h *APIHandler) listVerifications(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	recs, err := h.Store.List(r.Context(), limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, recs)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
