// Package takedown produces correctly formatted abuse-reports for
// registrars and hosting providers. The output is text the NGO
// lawyer reviews and sends — Joblantern never auto-submits.
package takedown

import (
	"bytes"
	"text/template"
	"time"
)

// Packet is the assembled report.
type Packet struct {
	Registrar       string
	Domain          string
	ReportedAt      time.Time
	Evidence        []Evidence
	ContactEmail    string
	ContactOrg      string
}

// Evidence is one item attached to the packet.
type Evidence struct {
	Kind        string
	Description string
	URL         string
	SHA256      string
}

const tpl = `Subject: Abuse report — {{.Domain}}
Date: {{.ReportedAt.Format "2006-01-02"}}

To: abuse@{{.Registrar}}
From: {{.ContactOrg}} <{{.ContactEmail}}>

We have evidence that {{.Domain}} is used to operate a recruitment-fraud scheme targeting migrant workers. Detail:

{{range .Evidence}}- {{.Kind}}: {{.Description}}
  URL:  {{.URL}}
  Hash: {{.SHA256}}
{{end}}
Please review under your acceptable-use policy. We can supply additional evidence on request.

{{.ContactOrg}}
`

// Render writes a registrar abuse packet.
func Render(p Packet) (string, error) {
	t, err := template.New("packet").Parse(tpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, p); err != nil {
		return "", err
	}
	return buf.String(), nil
}
