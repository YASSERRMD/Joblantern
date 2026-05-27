// Package privacy enforces that user identities never appear in the
// graph layer. Only entity-side nodes (phone, email, address) are
// retained; submitter ids are stripped at ingestion.
package privacy

// Allowed lists the entity kinds that may appear in published or
// exported graphs.
var Allowed = map[string]bool{
	"phone":    true,
	"email":    true,
	"address":  true,
	"director": true,
	"payment":  true,
}

// Scrub filters a slice of (kind,id) pairs down to the allowed kinds.
func Scrub(in []struct{ Kind, ID string }) []struct{ Kind, ID string } {
	out := in[:0]
	for _, x := range in {
		if Allowed[x.Kind] {
			out = append(out, x)
		}
	}
	return out
}
