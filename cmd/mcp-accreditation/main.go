// mcp-accreditation wraps national higher-education accreditation
// registries and exposes lookup tools to the agent.
package main

import (
	"log/slog"
	"os"
	"strings"
)

// Registry is one national accreditation registry. New registries can
// be added in deployment config without code changes.
type Registry struct {
	ID            string
	Country       string
	Name          string
	BaseURL       string
	NormalizeName func(string) string
}

// known is the seed list. Entries are documented and link only to
// public, queryable registries.
var known = []Registry{
	{ID: "us-chea", Country: "US", Name: "Council for Higher Education Accreditation", BaseURL: "https://www.chea.org"},
	{ID: "uk-ofs", Country: "UK", Name: "Office for Students", BaseURL: "https://www.officeforstudents.org.uk"},
	{ID: "de-akkr", Country: "DE", Name: "Akkreditierungsrat", BaseURL: "https://www.akkreditierungsrat.de"},
	{ID: "au-teqsa", Country: "AU", Name: "TEQSA", BaseURL: "https://www.teqsa.gov.au"},
	{ID: "in-ugc", Country: "IN", Name: "University Grants Commission", BaseURL: "https://www.ugc.gov.in"},
}

// canonical normalises an institution name for fuzzy matching against
// registry entries.
func canonical(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	// Pre-canonicalise registry names so lookups are O(1) once tools land.
	index := make(map[string]Registry, len(known))
	for _, r := range known {
		index[canonical(r.Name)] = r
	}
	log.Info("mcp-accreditation starting", "registries", len(index))
	// MCP server boilerplate follows the established pattern.
}
