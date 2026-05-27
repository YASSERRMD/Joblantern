// Package monitor watches for re-registration of takedown'd scam
// domains under variant TLDs.
package monitor

import "strings"

// Watch is one ongoing monitor.
type Watch struct {
	Base     string   // example.com
	TLDs     []string // .net, .org, .info, ...
	NotifyTo string
}

// Variants returns the candidate variant hostnames a respawned scam
// might use.
func (w Watch) Variants() []string {
	dot := strings.LastIndex(w.Base, ".")
	if dot <= 0 {
		return nil
	}
	root := w.Base[:dot]
	out := make([]string, 0, len(w.TLDs))
	for _, t := range w.TLDs {
		t = strings.TrimPrefix(t, ".")
		out = append(out, root+"."+t)
	}
	return out
}
