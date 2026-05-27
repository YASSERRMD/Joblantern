// Package thirdparty mirrors annual archives to durable third-party
// repositories — the Internet Archive and Zenodo — so the data
// survives even if Joblantern.org sunsets.
package thirdparty

// Target is one mirror destination.
type Target struct {
	ID   string
	Name string
	URL  string
}

// Defaults captures the seed list.
func Defaults() []Target {
	return []Target{
		{ID: "ia", Name: "Internet Archive", URL: "https://archive.org"},
		{ID: "zenodo", Name: "Zenodo (CERN)", URL: "https://zenodo.org"},
	}
}
