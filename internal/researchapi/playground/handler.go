// Package playground serves a minimal, self-hosted GraphQL explorer
// at /graphql/explorer. The HTML is a single page that posts to
// /graphql; no third-party CDN is required.
package playground

import "net/http"

const page = `<!doctype html>
<html><head><meta charset="utf-8"><title>Joblantern GraphQL</title>
<style>body{font-family:system-ui,sans-serif;margin:0;padding:1rem;background:#0e0e10;color:#e3e3e3}
textarea,pre{width:100%;box-sizing:border-box;background:#1a1a1d;color:#e3e3e3;border:1px solid #333;padding:.5rem;font-family:ui-monospace,monospace}
textarea{height:14rem}
button{background:#3b6;color:#000;border:0;padding:.5rem 1rem;cursor:pointer;margin-top:.5rem}
</style></head><body>
<h1>Joblantern GraphQL Explorer</h1>
<p>Read-only, anonymized verdict API. <a href="/docs/RESEARCH-API.md" style="color:#9cf">Docs</a></p>
<textarea id="q">{ verdictsByCountry { key count } }</textarea>
<button onclick="run()">Run</button>
<pre id="o"></pre>
<script>
async function run(){
  const r = await fetch('/graphql',{method:'POST',headers:{'content-type':'application/json'},body:JSON.stringify({query:document.getElementById('q').value})});
  document.getElementById('o').textContent = JSON.stringify(await r.json(), null, 2);
}
</script></body></html>`

// Handler returns the explorer page.
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(page))
	})
}
