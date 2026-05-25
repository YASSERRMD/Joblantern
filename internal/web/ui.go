package web

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/yasserrmd/joblantern/internal/agent"
)

//go:embed templates/*.html
var tmplFS embed.FS

// UI wires server-rendered pages on top of chi.
type UI struct {
	pages map[string]*template.Template
	store agent.Store
	api   *APIHandler
}

// pageData is the standard envelope every page receives.
type pageData struct {
	Title string
	Data  any
}

// NewUI registers routes "/" and "/verifications/:id".
func NewUI(r chi.Router, store agent.Store, api *APIHandler) (*UI, error) {
	funcs := template.FuncMap{
		"mul": func(a, b float64) float64 { return a * b },
	}
	pages := map[string]*template.Template{}
	for _, name := range []string{"home.html", "result.html"} {
		t, err := template.New(name).Funcs(funcs).ParseFS(tmplFS,
			"templates/layout.html", "templates/"+name)
		if err != nil {
			return nil, err
		}
		pages[name] = t
	}
	u := &UI{pages: pages, store: store, api: api}
	r.Get("/", u.home)
	r.Post("/verify", u.formSubmit)
	r.Get("/verifications/{id}", u.show)
	return u, nil
}

func (u *UI) home(w http.ResponseWriter, _ *http.Request) {
	u.render(w, "home.html", pageData{Title: "Joblantern – verify a job listing"})
}

// formSubmit accepts a posted HTML form, kicks off async verification,
// and redirects to the result page.
func (u *UI) formSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	sub := agent.Submission{
		ListingURL:     strings.TrimSpace(r.PostFormValue("listing_url")),
		ListingText:    r.PostFormValue("listing_text"),
		CompanyName:    strings.TrimSpace(r.PostFormValue("company_name")),
		ClaimedAddress: strings.TrimSpace(r.PostFormValue("claimed_address")),
		RecruiterEmail: strings.TrimSpace(r.PostFormValue("recruiter_email")),
		RecruiterPhone: strings.TrimSpace(r.PostFormValue("recruiter_phone")),
		Jurisdiction:   strings.ToUpper(strings.TrimSpace(r.PostFormValue("jurisdiction"))),
		Role:           strings.TrimSpace(r.PostFormValue("role")),
		Domain:         strings.TrimSpace(r.PostFormValue("domain")),
	}
	id, err := u.store.Create(r.Context(), sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	go u.api.RunAndStore(r.Context(), id, sub)
	http.Redirect(w, r, "/verifications/"+id, http.StatusSeeOther)
}

func (u *UI) show(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	rec, err := u.store.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rec == nil {
		http.NotFound(w, r)
		return
	}
	u.render(w, "result.html", pageData{Title: "Verification " + id, Data: rec})
}

func (u *UI) render(w http.ResponseWriter, name string, data pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := u.pages[name].ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
