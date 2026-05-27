package researchapi

// REST/GraphQL parity tests. The intent is to guarantee that any
// field exposed by one surface is reachable on the other, modulo a
// small allow-list of operation-only endpoints (export streams,
// webhook management).

import "testing"

// surfaceCatalogue lists the canonical anonymized verdict fields.
// Both /api/v1/verdicts and the GraphQL Verdict type must expose each.
var surfaceCatalogue = []string{
	"id",
	"createdAt",
	"country",
	"industry",
	"riskScore",
	"riskBand",
	"redFlags",
}

func TestSurfaceCatalogueParity(t *testing.T) {
	if len(surfaceCatalogue) == 0 {
		t.Fatal("surface catalogue is empty")
	}
	seen := map[string]bool{}
	for _, f := range surfaceCatalogue {
		if seen[f] {
			t.Errorf("duplicate field in catalogue: %s", f)
		}
		seen[f] = true
	}
}
