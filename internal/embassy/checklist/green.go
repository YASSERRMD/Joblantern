// Package checklist generates a green-verdict travel-readiness
// receipt. The intent is to surface the small, easy-to-forget items
// that even a legitimate offer needs in order.
package checklist

import "io"

// Item is one checklist line.
type Item struct {
	Code      string
	Label     string
	Mandatory bool
}

// Default returns the canonical pre-departure checklist.
func Default() []Item {
	return []Item{
		{"passport-validity", "Passport valid > 6 months from arrival", true},
		{"work-visa", "Work visa stamped (not tourist)", true},
		{"contract-copy", "Signed contract in your native language", true},
		{"agency-licence", "Agency licence number printed on contract", true},
		{"emergency-contacts", "Embassy hotline saved offline in phone", true},
		{"insurance", "Health insurance certificate carried", true},
		{"return-ticket", "Return ticket or repatriation clause confirmed", true},
		{"first-paycheck", "First paycheck schedule confirmed in writing", false},
	}
}

// Print writes the checklist to w.
func Print(w io.Writer, items []Item) error {
	for _, it := range items {
		box := "[ ]"
		mark := ""
		if it.Mandatory {
			mark = " *"
		}
		if _, err := io.WriteString(w, box+" "+it.Label+mark+"\n"); err != nil {
			return err
		}
	}
	return nil
}
