// Package domain combines three open data sources to profile a
// domain's age and history: port-43 WHOIS, crt.sh certificate
// transparency, and the Internet Archive Wayback Machine.
package domain

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

// WHOIS represents the few fields mcp-domain cares about.
type WHOIS struct {
	Domain       string    `json:"domain"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	Registrar    string    `json:"registrar,omitempty"`
	Registrant   string    `json:"registrant,omitempty"`
	CountryCode  string    `json:"country_code,omitempty"`
	Raw          string    `json:"raw,omitempty"`
	Unredacted   bool      `json:"unredacted,omitempty"`
}

// WHOISLookup is the interface implemented by both the real port-43
// resolver and a test stub.
type WHOISLookup interface {
	Lookup(ctx context.Context, domain string) (*WHOIS, error)
}

// PortLookup talks WHOIS over port 43 to a configurable server.
type PortLookup struct {
	Server  string
	Dialer  *net.Dialer
	Timeout time.Duration
}

// NewPortLookup returns a PortLookup pointed at whois.iana.org.
func NewPortLookup() *PortLookup {
	return &PortLookup{
		Server:  "whois.iana.org:43",
		Dialer:  &net.Dialer{Timeout: 10 * time.Second},
		Timeout: 15 * time.Second,
	}
}

// ErrInvalidDomain signals a malformed input domain.
var ErrInvalidDomain = errors.New("domain: invalid domain")

// Lookup performs a single-hop WHOIS query and parses the few
// human-readable fields we care about. The Raw response is preserved
// for diagnostics.
func (p *PortLookup) Lookup(ctx context.Context, dom string) (*WHOIS, error) {
	if !validDomain(dom) {
		return nil, ErrInvalidDomain
	}
	dctx, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()

	conn, err := p.Dialer.DialContext(dctx, "tcp", p.Server)
	if err != nil {
		return nil, fmt.Errorf("dial whois: %w", err)
	}
	defer func() { _ = conn.Close() }()

	if d, ok := dctx.Deadline(); ok {
		_ = conn.SetDeadline(d)
	}
	if _, err := fmt.Fprintf(conn, "%s\r\n", dom); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}

	var b strings.Builder
	sc := bufio.NewScanner(conn)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	for sc.Scan() {
		b.WriteString(sc.Text())
		b.WriteByte('\n')
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	w := parseWHOIS(dom, b.String())
	return w, nil
}

// validDomain is a very lightweight syntactic gate (not RFC-grade).
var validDomainRE = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)

func validDomain(d string) bool {
	return validDomainRE.MatchString(d)
}

var (
	createdRE = regexp.MustCompile(`(?i)(Creation Date|Created On|created|registered on|Domain Registration Date)\s*[:\.]+\s*([0-9T:\-\.+Zz/ ]+)`)
	registrarRE = regexp.MustCompile(`(?i)Registrar(?: Name)?\s*[:\.]+\s*(.+)`)
	countryRE = regexp.MustCompile(`(?i)Registrant Country\s*[:\.]+\s*([A-Za-z]{2,})`)
)

var dateFormats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05Z",
	"2006-01-02 15:04:05",
	"2006-01-02",
	"02-Jan-2006",
	"2006/01/02",
}

func parseWHOIS(dom, raw string) *WHOIS {
	w := &WHOIS{Domain: dom, Raw: truncate(raw, 4096)}
	if m := createdRE.FindStringSubmatch(raw); len(m) >= 3 {
		s := strings.TrimSpace(m[2])
		for _, f := range dateFormats {
			if t, err := time.Parse(f, s); err == nil {
				w.CreatedAt = t.UTC()
				break
			}
		}
	}
	if m := registrarRE.FindStringSubmatch(raw); len(m) >= 2 {
		w.Registrar = strings.TrimSpace(m[1])
	}
	if m := countryRE.FindStringSubmatch(raw); len(m) >= 2 {
		w.CountryCode = strings.ToUpper(strings.TrimSpace(m[1]))
	}
	// We don't try to parse the Registrant name — usually redacted for
	// GDPR. Mark Unredacted=false unless we discovered a country (a
	// rough proxy).
	w.Unredacted = w.CountryCode != ""
	return w
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
