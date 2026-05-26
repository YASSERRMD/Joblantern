package web

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/yasserrmd/joblantern/internal/transparency"
)

// VerdictSource lets the dashboard pull anonymised rollups without
// touching agent / db internals directly. The joblantern binary wires
// it to a snapshot loaded from agent.Store; an alternative impl can
// read from a materialised view once the nightly cron lands.
type VerdictSource func() []transparency.Verdict

// TransparencyHandler exposes:
//
//	GET /transparency               minimal HTML landing
//	GET /transparency/aggregate     JSON rows
//	GET /transparency/aggregate.csv CSV download
//	GET /transparency/aggregate.jsonl JSONL one-row-per-line export
type TransparencyHandler struct {
	source VerdictSource
	agg    *transparency.Aggregator

	mu       sync.Mutex
	cached   []transparency.Row
	cachedAt time.Time
	cacheTTL time.Duration
}

// NewTransparency registers routes on r.
func NewTransparency(r chi.Router, src VerdictSource) *TransparencyHandler {
	h := &TransparencyHandler{
		source:   src,
		agg:      transparency.New(),
		cacheTTL: 5 * time.Minute,
	}
	r.Get("/transparency", h.landing)
	r.Get("/transparency/aggregate", h.json)
	r.Get("/transparency/aggregate.csv", h.csv)
	r.Get("/transparency/aggregate.jsonl", h.jsonl)
	return h
}

func (h *TransparencyHandler) rows() []transparency.Row {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.cached != nil && time.Since(h.cachedAt) < h.cacheTTL {
		return h.cached
	}
	in := h.source()
	h.cached = h.agg.Aggregate(in)
	h.cachedAt = time.Now()
	return h.cached
}

func (h *TransparencyHandler) json(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(h.rows())
}

func (h *TransparencyHandler) csv(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"joblantern-aggregate.csv\"")
	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"date", "country", "risk", "count", "fuzzed"})
	for _, r := range h.rows() {
		_ = cw.Write([]string{
			r.Date, r.Country, r.Risk,
			strconv.Itoa(r.Count),
			strconv.FormatBool(r.Fuzzed),
		})
	}
	cw.Flush()
}

func (h *TransparencyHandler) jsonl(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Content-Disposition", "attachment; filename=\"joblantern-aggregate.jsonl\"")
	enc := json.NewEncoder(w)
	for _, r := range h.rows() {
		_ = enc.Encode(r)
	}
}

func (h *TransparencyHandler) landing(w http.ResponseWriter, _ *http.Request) {
	rows := h.rows()
	// Roll up totals for the landing summary.
	totals := map[string]int{}
	countries := map[string]struct{}{}
	for _, r := range rows {
		totals[r.Risk] += r.Count
		if r.Country != "" {
			countries[r.Country] = struct{}{}
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en"><head>
<meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<title>Joblantern — transparency</title>
<style>
body{font-family:-apple-system,system-ui,sans-serif;background:#0b0d12;color:#e8eaef;margin:0;padding:2rem}
main{max-width:760px;margin:0 auto}
h1{font-weight:600;margin:0 0 .5rem}
.muted{color:#8b8f9a}
.cards{display:flex;gap:1rem;margin-top:1rem}
.card{flex:1;background:#11141c;border:1px solid #2a2f3d;border-radius:8px;padding:1rem}
.card .n{font-size:1.6rem;font-weight:600}
.green{color:#3da76a}.yellow{color:#d6a23a}.red{color:#d6452a}
a{color:#7aa2ff}
table{width:100%%;margin-top:1rem;border-collapse:collapse}
th,td{padding:.4rem .6rem;border-bottom:1px solid #2a2f3d;text-align:left;font-size:.92rem}
</style>
</head><body><main>
<h1>Transparency</h1>
<p class="muted">Anonymised verdict aggregates. Small cells dropped, counts fuzzed with discrete Laplace noise. No individual user is identifiable.</p>

<div class="cards">
  <div class="card"><div class="muted">Green</div><div class="n green">%d</div></div>
  <div class="card"><div class="muted">Yellow</div><div class="n yellow">%d</div></div>
  <div class="card"><div class="muted">Red</div><div class="n red">%d</div></div>
  <div class="card"><div class="muted">Countries</div><div class="n">%d</div></div>
</div>

<p style="margin-top:1.5rem">
  Download:
  <a href="/transparency/aggregate">JSON</a> ·
  <a href="/transparency/aggregate.csv">CSV</a> ·
  <a href="/transparency/aggregate.jsonl">JSONL</a>
</p>

<h2>Recent rows</h2>
<table><thead><tr><th>Date</th><th>Country</th><th>Risk</th><th>Count</th></tr></thead><tbody>`,
		totals["green"], totals["yellow"], totals["red"], len(countries))

	maxRows := 50
	for i, r := range rows {
		if i >= maxRows {
			break
		}
		fmt.Fprintf(w, `<tr><td>%s</td><td>%s</td><td class="%s">%s</td><td>%d</td></tr>`,
			r.Date, r.Country, r.Risk, r.Risk, r.Count)
	}
	fmt.Fprint(w, `</tbody></table>
<p class="muted" style="margin-top:2rem;font-size:.8rem">
Joblantern · Apache 2.0 · <a href="/">Home</a>
</p>
</main></body></html>`)
}
