package web

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed static
var staticFS embed.FS

// MountStatic registers /static/* plus root aliases for sw.js and the
// manifest, so the service worker's scope can be the application root.
func MountStatic(r chi.Router) error {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		return err
	}
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(sub))))

	r.Get("/sw.js", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("Service-Worker-Allowed", "/")
		http.ServeFileFS(w, req, sub, "sw.js")
	})
	r.Get("/manifest.webmanifest", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/manifest+json")
		http.ServeFileFS(w, req, sub, "manifest.webmanifest")
	})
	return nil
}
