// Package landing renders the citable archive landing page. Each
// archive has a stable URL that returns either HTML for browsers or
// BibTeX for citation managers.
package landing

import "fmt"

// Citation is the metadata that drives the landing page.
type Citation struct {
	DOI     string
	Year    int
	URL     string
	Title   string
	Bytes   int64
	Records int
	SHA256  string
}

// BibTeX returns the BibTeX entry.
func (c Citation) BibTeX() string {
	return fmt.Sprintf(`@misc{joblantern%d,
  title = {%s},
  author = {Joblantern Foundation},
  year = {%d},
  publisher = {Joblantern Foundation},
  doi = {%s},
  url = {%s}
}`, c.Year, c.Title, c.Year, c.DOI, c.URL)
}

// HTML returns the landing page body (no external assets).
func (c Citation) HTML() string {
	return fmt.Sprintf(`<!doctype html><html><head><meta charset="utf-8"><title>%s</title></head>
<body>
<h1>%s</h1>
<p>DOI: <code>%s</code></p>
<p>SHA-256: <code>%s</code></p>
<p>Records: %d. Bytes: %d.</p>
<pre>%s</pre>
<p><a href="%s">Download</a></p>
</body></html>`, c.Title, c.Title, c.DOI, c.SHA256, c.Records, c.Bytes, c.BibTeX(), c.URL)
}
