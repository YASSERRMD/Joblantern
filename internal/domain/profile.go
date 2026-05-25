package domain

import (
	"context"
	"time"
)

// Profile composes WHOIS, SSL history, and archive history into a
// single record. FreshnessScore approaches 1.0 when the domain looks
// brand new (and therefore suspicious for a company claiming long
// operation), and approaches 0.0 for well-established domains.
type Profile struct {
	Domain         string          `json:"domain"`
	WHOIS          *WHOIS          `json:"whois,omitempty"`
	SSL            *CertSummary    `json:"ssl,omitempty"`
	Archive        *ArchiveSummary `json:"archive,omitempty"`
	AgeDays        int             `json:"age_days"`
	FreshnessScore float64         `json:"freshness_score"`
}

// Composer pulls all three sources in parallel and computes a Profile.
type Composer struct {
	WHOIS   WHOISLookup
	CrtSH   *CrtSHClient
	Wayback *WaybackClient
}

// NewComposer wires the standard backends.
func NewComposer() *Composer {
	return &Composer{
		WHOIS:   NewPortLookup(),
		CrtSH:   NewCrtSHClient(),
		Wayback: NewWaybackClient(),
	}
}

// FullProfile runs all three lookups. Failures in individual backends
// are recorded as nil fields; the overall call only errors if all
// three failed.
func (c *Composer) FullProfile(ctx context.Context, dom string) (*Profile, error) {
	p := &Profile{Domain: dom, AgeDays: -1}

	type whoisRes struct {
		w   *WHOIS
		err error
	}
	type sslRes struct {
		s   *CertSummary
		err error
	}
	type archiveRes struct {
		a   *ArchiveSummary
		err error
	}

	wch := make(chan whoisRes, 1)
	sch := make(chan sslRes, 1)
	ach := make(chan archiveRes, 1)

	go func() {
		w, err := c.WHOIS.Lookup(ctx, dom)
		wch <- whoisRes{w, err}
	}()
	go func() {
		s, err := c.CrtSH.Summary(ctx, dom)
		sch <- sslRes{s, err}
	}()
	go func() {
		a, err := c.Wayback.Summary(ctx, dom)
		ach <- archiveRes{a, err}
	}()

	w := <-wch
	s := <-sch
	a := <-ach
	if w.err == nil {
		p.WHOIS = w.w
	}
	if s.err == nil {
		p.SSL = s.s
	}
	if a.err == nil {
		p.Archive = a.a
	}
	if p.WHOIS == nil && p.SSL == nil && p.Archive == nil {
		// All failed; surface the first error we have.
		if w.err != nil {
			return p, w.err
		}
		if s.err != nil {
			return p, s.err
		}
		return p, a.err
	}

	// Derive age.
	var oldest time.Time
	if p.WHOIS != nil && !p.WHOIS.CreatedAt.IsZero() {
		oldest = p.WHOIS.CreatedAt
	}
	if p.SSL != nil && !p.SSL.FirstCertAt.IsZero() {
		if oldest.IsZero() || p.SSL.FirstCertAt.Before(oldest) {
			oldest = p.SSL.FirstCertAt
		}
	}
	if p.Archive != nil && !p.Archive.EarliestSnapshot.IsZero() {
		if oldest.IsZero() || p.Archive.EarliestSnapshot.Before(oldest) {
			oldest = p.Archive.EarliestSnapshot
		}
	}
	if !oldest.IsZero() {
		p.AgeDays = int(time.Since(oldest).Hours() / 24)
	}

	// Freshness score: brand new = 1.0, ~5 years old or more = ~0.0.
	switch {
	case p.AgeDays < 0:
		p.FreshnessScore = 0.5
	case p.AgeDays < 90:
		p.FreshnessScore = 1.0
	case p.AgeDays > 365*5:
		p.FreshnessScore = 0.0
	default:
		p.FreshnessScore = 1.0 - float64(p.AgeDays)/float64(365*5)
	}
	return p, nil
}
