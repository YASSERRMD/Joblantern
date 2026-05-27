// Package ui serves the kiosk HTML — large text, no keyboard input
// required, scan-driven. All localisation is loaded from the existing
// i18n bundle from Phase 22.
package ui

import (
	"html/template"
	"net/http"
)

const kioskHTML = `<!doctype html>
<html lang="{{.Lang}}"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Joblantern Kiosk</title>
<style>
:root{font-size:22pt}
body{font-family:system-ui,sans-serif;margin:0;padding:2rem;background:#0b1220;color:#fff}
.btn{display:block;width:100%;padding:2rem;margin:1rem 0;background:#3b6;color:#000;font-size:1.4rem;border:0;border-radius:1rem;cursor:pointer}
.hint{color:#9cf;margin-top:1rem}
.scan{font-size:1.6rem;text-align:center;padding:3rem;border:4px dashed #3b6;border-radius:1rem;margin-top:2rem}
</style></head><body>
<h1>{{.WelcomeTitle}}</h1>
<p class="hint">{{.WelcomeSub}}</p>
<div class="scan">{{.ScanPrompt}}</div>
<button class="btn" onclick="location='/kiosk/officer'">{{.OfficerLabel}}</button>
</body></html>`

// Page is the model passed to the kiosk template.
type Page struct {
	Lang         string
	WelcomeTitle string
	WelcomeSub   string
	ScanPrompt   string
	OfficerLabel string
}

// Handler serves the kiosk landing page.
func Handler(t *template.Template, page Page) http.Handler {
	if t == nil {
		t = template.Must(template.New("kiosk").Parse(kioskHTML))
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = t.Execute(w, page)
	})
}
