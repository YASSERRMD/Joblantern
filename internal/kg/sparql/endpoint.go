// Package sparql is the HTTP endpoint that serves SPARQL 1.1 queries.
// Cost limits mirror the GraphQL cost analyser so the surface is hard
// to abuse.
package sparql

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Cap is the maximum result row count for any query.
const Cap = 5000

// Endpoint is the HTTP handler.
type Endpoint struct {
	Execute func(query string) ([]map[string]string, error)
}

// ServeHTTP implements http.Handler.
func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query().Get("query")
	if q == "" && r.Method == http.MethodPost {
		_ = r.ParseForm()
		q = r.PostForm.Get("query")
	}
	if q == "" {
		http.Error(w, "query required", http.StatusBadRequest)
		return
	}
	rows, err := e.Execute(q)
	if err != nil {
		if errors.Is(err, errCapped) {
			http.Error(w, "result capped at "+itoa(Cap), http.StatusBadRequest)
			return
		}
		http.Error(w, "query failed", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/sparql-results+json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"head":    map[string]any{"vars": varsOf(rows)},
		"results": map[string]any{"bindings": rows},
	})
}

var errCapped = errors.New("result capped")

func varsOf(rows []map[string]string) []string {
	if len(rows) == 0 {
		return nil
	}
	out := make([]string, 0, len(rows[0]))
	for k := range rows[0] {
		out = append(out, k)
	}
	return out
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [12]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
