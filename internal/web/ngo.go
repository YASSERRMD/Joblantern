package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jung-kurt/gofpdf"

	"github.com/yasserrmd/joblantern/internal/agent"
)

// NGOHandler bundles two small NGO-facing surfaces:
//
//	GET /kiosk                  large-text walk-in form (no chrome, no JS)
//	POST /kiosk/verify          same as /verify but redirects to /kiosk/result/:id
//	GET /kiosk/result/{id}      large-text result page (auto-refresh)
//	GET /verifications/{id}/print.pdf  single-page A4 PDF in user's lang
//
// Kiosk mode is intentionally minimal — designed for in-person clinics
// where a caseworker turns a tablet around for a walk-in client.
type NGOHandler struct {
	store agent.Store
	api   *APIHandler
}

// NewNGO registers the routes on r.
func NewNGO(r chi.Router, store agent.Store, api *APIHandler) *NGOHandler {
	h := &NGOHandler{store: store, api: api}
	r.Get("/kiosk", h.kioskHome)
	r.Post("/kiosk/verify", h.kioskSubmit)
	r.Get("/kiosk/result/{id}", h.kioskResult)
	r.Get("/verifications/{id}/print.pdf", h.printPDF)
	return h
}

func (h *NGOHandler) kioskHome(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Joblantern · kiosk</title>
<style>
body{font-family:-apple-system,system-ui,sans-serif;background:#fff;color:#111;margin:0;padding:2rem;font-size:20px}
h1{font-size:2rem;margin:0 0 .5rem}
label{display:block;font-size:1.2rem;margin-top:1.2rem}
textarea,input,select{width:100%;font-size:1.1rem;padding:.75rem;border:2px solid #444;border-radius:8px;box-sizing:border-box}
textarea{min-height:14rem}
button{margin-top:1.5rem;padding:1rem 2rem;font-size:1.2rem;background:#1a4ea0;color:#fff;border:0;border-radius:8px;cursor:pointer}
main{max-width:680px;margin:0 auto}
</style></head><body><main>
<h1>Is this job real?</h1>
<p>Paste the job offer below. A caseworker will review the verdict with you.</p>
<form action="/kiosk/verify" method="post">
  <label>Paste the message or offer text
    <textarea name="listing_text" required></textarea>
  </label>
  <label>Company name (if you know it)
    <input name="company_name">
  </label>
  <label>Destination country (e.g. AE, SA, BD)
    <input name="jurisdiction" maxlength="2">
  </label>
  <button type="submit">Check</button>
</form>
</main></body></html>`))
}

func (h *NGOHandler) kioskSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	sub := agent.Submission{
		ListingText:  r.PostFormValue("listing_text"),
		CompanyName:  strings.TrimSpace(r.PostFormValue("company_name")),
		Jurisdiction: strings.ToUpper(strings.TrimSpace(r.PostFormValue("jurisdiction"))),
	}
	id, err := h.store.Create(r.Context(), sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	go h.api.RunAndStore(r.Context(), id, sub)
	http.Redirect(w, r, "/kiosk/result/"+id, http.StatusSeeOther)
}

func (h *NGOHandler) kioskResult(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	rec, _ := h.store.Get(r.Context(), id)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if rec == nil {
		http.NotFound(w, r)
		return
	}
	risk, color, headline := "pending", "#888", "Working…"
	if rec.Verdict != nil {
		risk = rec.Verdict.OverallRisk
		switch risk {
		case "red":
			color = "#d6452a"
			headline = "This looks unsafe."
		case "yellow":
			color = "#d6a23a"
			headline = "Be careful — verify in person."
		case "green":
			color = "#3da76a"
			headline = "No major warnings found."
		}
	}
	refresh := ""
	if rec.Status != "completed" && rec.Status != "failed" {
		refresh = `<meta http-equiv="refresh" content="3">`
	}
	reasons := ""
	if rec.Verdict != nil {
		for i, rs := range rec.Verdict.Reasons {
			if i >= 5 {
				break
			}
			reasons += "<li>" + esc(rs) + "</li>"
		}
	}
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en"><head>%s
<meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Result · Joblantern</title>
<style>
body{font-family:-apple-system,system-ui,sans-serif;background:#fff;color:#111;margin:0;padding:2rem;font-size:22px}
main{max-width:680px;margin:0 auto;text-align:center}
.badge{display:inline-block;padding:1rem 2rem;border-radius:14px;background:%s;color:#fff;font-weight:700;font-size:1.8rem;margin:1rem 0}
h1{font-size:2rem}
ul{text-align:left;font-size:1.1rem}
a{color:#1a4ea0}
</style></head><body><main>
<h1>%s</h1>
<div class="badge">%s</div>
<p>Verification id: <code>%s</code></p>
%s
<p>
  <a href="/verifications/%s/print.pdf">Print one-page summary</a> ·
  <a href="/kiosk">Check another</a>
</p>
</main></body></html>`,
		refresh, color, esc(headline), strings.ToUpper(risk), id,
		ifEmpty(reasons, "", "<h2>Top reasons</h2><ul>"+reasons+"</ul>"),
		id)
}

func (h *NGOHandler) printPDF(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	rec, _ := h.store.Get(r.Context(), id)
	if rec == nil {
		http.NotFound(w, r)
		return
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Helvetica", "B", 22)
	pdf.Cell(0, 12, "Joblantern verification")
	pdf.Ln(14)

	pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(0, 6, "Generated "+time.Now().UTC().Format("2006-01-02 15:04 UTC"))
	pdf.Ln(6)
	pdf.Cell(0, 6, "ID: "+rec.ID)
	pdf.Ln(10)

	// Verdict
	if rec.Verdict != nil {
		pdf.SetFont("Helvetica", "B", 28)
		switch rec.Verdict.OverallRisk {
		case "red":
			pdf.SetTextColor(214, 69, 42)
		case "yellow":
			pdf.SetTextColor(214, 162, 58)
		default:
			pdf.SetTextColor(61, 167, 106)
		}
		pdf.Cell(0, 14, strings.ToUpper(rec.Verdict.OverallRisk))
		pdf.SetTextColor(0, 0, 0)
		pdf.Ln(16)
		pdf.SetFont("Helvetica", "", 11)
		pdf.Cell(0, 6, fmt.Sprintf("Confidence: %.0f%%", rec.Verdict.Confidence*100))
		pdf.Ln(10)

		if len(rec.Verdict.Reasons) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.Cell(0, 6, "Reasons")
			pdf.Ln(7)
			pdf.SetFont("Helvetica", "", 11)
			for i, rs := range rec.Verdict.Reasons {
				if i >= 8 {
					break
				}
				pdf.MultiCell(0, 6, "  • "+rs, "", "", false)
			}
			pdf.Ln(4)
		}
	} else {
		pdf.SetFont("Helvetica", "", 12)
		pdf.MultiCell(0, 6, "Verdict not yet available. Status: "+rec.Status, "", "", false)
	}

	// Footer
	pdf.SetY(-25)
	pdf.SetFont("Helvetica", "I", 9)
	pdf.SetTextColor(120, 120, 120)
	pdf.MultiCell(0, 4,
		"Joblantern is not a lawyer. Use this verdict as one input among many; "+
			"always verify directly with the destination country's labour regulator.",
		"", "", false)

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="joblantern-%s.pdf"`, id))
	if err := pdf.Output(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func esc(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&#39;").Replace(s)
}

func ifEmpty(test, ifEmpty, ifSet string) string {
	if test == "" {
		return ifEmpty
	}
	return ifSet
}
