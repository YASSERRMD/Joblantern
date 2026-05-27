// Package partnerships tracks the consumer-protection-agency
// partnerships specific to the marketplace module. These are typically
// not the labour ministries from Phase 37 but parallel consumer
// agencies.
package partnerships

// Agency is one consumer-protection partner.
type Agency struct {
	ID         string
	Country    string
	Name       string
	ContactURL string
	MOUSigned  bool
}

// Defaults returns the seed list.
func Defaults() []Agency {
	return []Agency{
		{Country: "PH", Name: "DTI Consumer Protection"},
		{Country: "IN", Name: "Consumer Affairs (MoCA)"},
		{Country: "AE", Name: "UAE Consumer Protection Department"},
		{Country: "DE", Name: "Verbraucherzentrale Bundesverband"},
	}
}
