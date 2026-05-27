// Package doi mints DOIs for each annual archive via DataCite.
package doi

import "fmt"

// Mint is the registered DOI plus the metadata sent to DataCite.
type Mint struct {
	DOI       string // "10.<prefix>/joblantern.<year>"
	Year      int
	URL       string
	Title     string
	Creator   string
	Publisher string
}

// Build constructs a Mint with a deterministic suffix.
func Build(prefix string, year int, url string) Mint {
	return Mint{
		DOI:       fmt.Sprintf("%s/joblantern.%d", prefix, year),
		Year:      year,
		URL:       url,
		Title:     fmt.Sprintf("Joblantern Annual Verdict Archive, %d", year),
		Creator:   "Joblantern Foundation",
		Publisher: "Joblantern Foundation",
	}
}
