// Package lookup queries mcp-scam-db.search_reports_by_phone over a
// stateless RPC and returns a normalised verdict suitable for an
// overlay.
package lookup

// Verdict is the trimmed response sent to the device.
type Verdict struct {
	Band    string `json:"band"`   // green / yellow / red
	Reason  string `json:"reason"` // single-line user-facing message
	Updated string `json:"updated_at,omitempty"`
}

// Classifier wraps the MCP scam-db client.
type Classifier interface {
	ClassifyHash(phoneHash string) (Verdict, error)
}

// Service holds dependencies.
type Service struct {
	C Classifier
}

// Lookup returns a verdict for a raw phone string. The number itself
// is hashed locally and only the hash crosses the wire — see
// internal/callerid/privacy for the hashing scheme.
func (s Service) Lookup(rawPhone string) Verdict {
	h := HashPhone(rawPhone)
	v, err := s.C.ClassifyHash(h)
	if err != nil {
		return Verdict{Band: "yellow", Reason: "Could not check — treat with caution"}
	}
	return v
}
