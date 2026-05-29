// Package sdk generates lightweight Python and R client packages from
// the OpenAPI document. The intent is not feature parity with
// hand-written SDKs — researchers should be able to `pip install` or
// `install.packages` and reach the API in one line.
package sdk

import (
	"fmt"
	"strings"
)

// Lang identifies a target language.
type Lang string

const (
	Python Lang = "python"
	R      Lang = "r"
)

// Generate produces source for a thin client. The OpenAPI spec is
// passed in raw form; the generator only needs the base URL and tier.
func Generate(lang Lang, baseURL string) (string, error) {
	switch lang {
	case Python:
		return python(baseURL), nil
	case R:
		return rLang(baseURL), nil
	}
	return "", fmt.Errorf("unsupported language: %s", lang)
}

func python(baseURL string) string {
	var b strings.Builder
	b.WriteString("\"\"\"Joblantern research SDK (auto-generated).\"\"\"\n")
	b.WriteString("import urllib.request, json\n\n")
	fmt.Fprintf(&b, "BASE = %q\n\n", baseURL)
	b.WriteString("class Client:\n")
	b.WriteString("    def __init__(self, token=None):\n        self.token = token\n\n")
	b.WriteString("    def _get(self, path, params=None):\n")
	b.WriteString("        url = BASE + path\n")
	b.WriteString("        req = urllib.request.Request(url)\n")
	b.WriteString("        if self.token: req.add_header('Authorization', f'Bearer {self.token}')\n")
	b.WriteString("        with urllib.request.urlopen(req) as r: return json.load(r)\n\n")
	b.WriteString("    def verdicts(self, country=None, since=None, cursor=None):\n")
	b.WriteString("        q = []\n        if country: q.append(f'country={country}')\n        if since: q.append(f'since={since}')\n        if cursor: q.append(f'cursor={cursor}')\n")
	b.WriteString("        return self._get('/api/v1/verdicts' + ('?' + '&'.join(q) if q else ''))\n")
	return b.String()
}

func rLang(baseURL string) string {
	var b strings.Builder
	b.WriteString("# Joblantern research SDK (auto-generated)\n")
	b.WriteString("library(httr); library(jsonlite)\n\n")
	fmt.Fprintf(&b, "joblantern_base <- %q\n\n", baseURL)
	b.WriteString("joblantern_verdicts <- function(token=NULL, country=NULL, since=NULL, cursor=NULL) {\n")
	b.WriteString("  q <- list(country=country, since=since, cursor=cursor)\n")
	b.WriteString("  q <- q[!sapply(q, is.null)]\n")
	b.WriteString("  url <- paste0(joblantern_base, '/api/v1/verdicts')\n")
	b.WriteString("  h <- if (!is.null(token)) add_headers(Authorization=paste('Bearer', token)) else NULL\n")
	b.WriteString("  r <- GET(url, query=q, h)\n")
	b.WriteString("  fromJSON(content(r, 'text', encoding='UTF-8'))\n")
	b.WriteString("}\n")
	return b.String()
}
