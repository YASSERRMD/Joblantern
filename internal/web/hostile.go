package web

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/yasserrmd/joblantern/internal/agent"
)

// HostileHandler wires the small set of endpoints aimed at operators
// in surveilled / unreliable / hostile-network conditions:
//
//	POST /panic-wipe          → clear all local browser + cookie state
//	GET  /lite                → minimum-bandwidth HTML home (no JS, no CSS hosting)
//
// The /panic-wipe response sets an expired cookie for every cookie
// the browser may have stored under our domain, and instructs the
// browser to clear its site data via Clear-Site-Data. We do not
// reach into IndexedDB from the server side; the response page
// includes a tiny JS snippet that wipes IndexedDB + caches.
type HostileHandler struct{}

// NewHostile registers the routes on r.
func NewHostile(r chi.Router, _ agent.Store) *HostileHandler {
	h := &HostileHandler{}
	r.Get("/panic-wipe", h.panic)
	r.Post("/panic-wipe", h.panic)
	r.Get("/lite", h.lite)
	return h
}

func (h *HostileHandler) panic(w http.ResponseWriter, _ *http.Request) {
	// Tell the browser to nuke as much locally-stored state as it can.
	// Clear-Site-Data is browser-enforced and respected by modern
	// Chromium / Firefox / Safari.
	w.Header().Set("Clear-Site-Data", `"cache","cookies","storage","executionContexts"`)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write([]byte(`<!DOCTYPE html>
<html lang="en"><head>
<meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<title>Joblantern — wiped</title>
<style>body{font-family:-apple-system,system-ui,sans-serif;background:#0b0d12;color:#e8eaef;display:flex;align-items:center;justify-content:center;min-height:100vh;margin:0;text-align:center;padding:1rem}main{max-width:380px}a{color:#7aa2ff}</style>
</head><body><main>
<h1>Cleared</h1>
<p>Browser caches, cookies, IndexedDB, and service-worker registrations have been removed for this site.</p>
<p><a href="/">Return home</a></p>
<script>
(async function(){
  try {
    if ('serviceWorker' in navigator) {
      var regs = await navigator.serviceWorker.getRegistrations();
      for (var r of regs) { try { await r.unregister(); } catch(_){} }
    }
    if ('caches' in window) {
      var keys = await caches.keys();
      for (var k of keys) { try { await caches.delete(k); } catch(_){} }
    }
    if (window.indexedDB && indexedDB.databases) {
      var dbs = await indexedDB.databases();
      for (var d of dbs) { try { indexedDB.deleteDatabase(d.name); } catch(_){} }
    } else {
      try { indexedDB.deleteDatabase("joblantern"); } catch(_){}
      try { indexedDB.deleteDatabase("joblantern-pwa"); } catch(_){}
    }
    try { localStorage.clear(); } catch(_){}
    try { sessionStorage.clear(); } catch(_){}
  } catch(e) {}
})();
</script>
</main></body></html>`))
}

func (h *HostileHandler) lite(w http.ResponseWriter, _ *http.Request) {
	// Plain HTML, no CSS, no JS. ~1 KB on the wire.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write([]byte(`<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8"><title>Joblantern (lite)</title></head>
<body>
<h1>Joblantern — lite mode</h1>
<p>Minimum-bandwidth form. Verdicts are text-only, no map tiles.</p>
<form action="/verify" method="post">
  <p>Paste the recruiter message:<br>
    <textarea name="listing_text" rows="10" cols="60"></textarea></p>
  <p>Company: <input name="company_name"></p>
  <p>Country (ISO-2): <input name="jurisdiction" maxlength="2"></p>
  <p><button type="submit">Verify</button></p>
</form>
<p><a href="/panic-wipe">Panic wipe</a> · <a href="/">Standard UI</a></p>
</body></html>`))
}

// Quickly avoid unused-import warning when the package is built without
// the agent dep.
var _ = time.Second
var _ = context.Background
