package graphql

import (
	"encoding/json"
	"net/http"
)

// Request is the wire form of a GraphQL request.
type Request struct {
	Query         string         `json:"query"`
	OperationName string         `json:"operationName,omitempty"`
	Variables     map[string]any `json:"variables,omitempty"`
}

// Handler is the HTTP handler for the /graphql endpoint. Implementers
// inject a Resolver (gqlgen-generated) and apply the complexity/cost
// middleware before this handler runs.
type Handler struct {
	Execute func(r Request) (any, []error)
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	data, errs := h.Execute(req)
	out := map[string]any{"data": data}
	if len(errs) > 0 {
		msgs := make([]map[string]any, len(errs))
		for i, e := range errs {
			msgs[i] = map[string]any{"message": e.Error()}
		}
		out["errors"] = msgs
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
