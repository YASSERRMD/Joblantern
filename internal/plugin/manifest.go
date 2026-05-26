// Package plugin defines Joblantern's third-party MCP plugin format.
//
// A plugin is described by a single YAML file (`joblantern-mcp.yaml`).
// Operators register plugins by URL or local path; the registry stores
// the manifest plus a trust label. Untrusted plugins are sandboxed
// (recommendation: gVisor or Firecracker) by the operator's container
// runtime — this package only models the manifest and registry, not
// the sandbox itself.
package plugin

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"time"

	"gopkg.in/yaml.v3"
)

// Trust is the operator's stance on a plugin.
type Trust string

const (
	TrustOfficial  Trust = "official"  // Maintained by the Joblantern project
	TrustCommunity Trust = "community" // Reviewed by operator; default trust
	TrustExternal  Trust = "external"  // Unreviewed; sandbox required
)

// Manifest is the on-disk representation of a plugin.
type Manifest struct {
	// Identity
	Name        string `yaml:"name"` // "joblantern.osint_people"
	DisplayName string `yaml:"display_name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`

	// Author + provenance
	Author    string `yaml:"author"`
	Homepage  string `yaml:"homepage"`
	License   string `yaml:"license"` // SPDX id; Joblantern accepts permissive only
	SourceURL string `yaml:"source_url"`

	// Runtime
	Transport   string            `yaml:"transport"` // stdio | http
	Command     string            `yaml:"command"`   // for stdio
	Args        []string          `yaml:"args"`
	URL         string            `yaml:"url"` // for http
	Env         map[string]string `yaml:"env"`
	CallTimeout time.Duration     `yaml:"call_timeout"`

	// Tools declared (informational only; the agent inspects the
	// server at runtime). Helps reviewers spot scope creep.
	Tools []ToolDecl `yaml:"tools"`

	// Attribution / share-alike. Must be displayed wherever the
	// plugin's output is shown.
	Attribution Attribution `yaml:"attribution"`

	// Signing — sigstore bundle path or detached signature; either
	// SignatureB64 + PubKeyHex must verify the canonical YAML, or
	// the plugin must be marked TrustExternal.
	SignatureB64 string `yaml:"signature_b64,omitempty"`
	PubKeyHex    string `yaml:"pubkey_hex,omitempty"`
}

// ToolDecl is the informational tool entry.
type ToolDecl struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// Attribution carries the required UI string + URL.
type Attribution struct {
	Text string `yaml:"text"`
	URL  string `yaml:"url"`
}

// LoadManifest parses a YAML document.
func LoadManifest(data []byte) (*Manifest, error) {
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("plugin: parse manifest: %w", err)
	}
	if err := m.validate(); err != nil {
		return nil, err
	}
	return &m, nil
}

// Permissive lists the licenses Joblantern will surface without warning.
var Permissive = map[string]bool{
	"Apache-2.0":   true,
	"MIT":          true,
	"BSD-2-Clause": true,
	"BSD-3-Clause": true,
	"ISC":          true,
	"MPL-2.0":      true,
}

func (m *Manifest) validate() error {
	if m.Name == "" {
		return errors.New("plugin: name is required")
	}
	if m.Version == "" {
		return errors.New("plugin: version is required")
	}
	if m.License == "" {
		return errors.New("plugin: license (SPDX id) is required")
	}
	if !Permissive[m.License] {
		return fmt.Errorf("plugin: license %q is not on the permissive allowlist", m.License)
	}
	switch m.Transport {
	case "stdio":
		if m.Command == "" {
			return errors.New("plugin: stdio transport requires command")
		}
	case "http":
		if m.URL == "" {
			return errors.New("plugin: http transport requires url")
		}
		if _, err := url.Parse(m.URL); err != nil {
			return fmt.Errorf("plugin: invalid url: %w", err)
		}
	default:
		return fmt.Errorf("plugin: transport must be stdio or http")
	}
	if m.Attribution.Text == "" {
		return errors.New("plugin: attribution.text is required")
	}
	return nil
}

// Verify checks the signature on the canonical YAML. Returns nil when
// either (a) signature verifies, or (b) the manifest carries neither
// a signature nor a pubkey and the trust label is "external".
func (m *Manifest) Verify(trust Trust) error {
	if m.SignatureB64 == "" || m.PubKeyHex == "" {
		if trust == TrustExternal {
			return nil
		}
		return errors.New("plugin: missing signature; only external trust may be unsigned")
	}
	pub, err := hexDecode(m.PubKeyHex)
	if err != nil {
		return fmt.Errorf("plugin: decode pubkey: %w", err)
	}
	if len(pub) != ed25519.PublicKeySize {
		return errors.New("plugin: pubkey must be 32 bytes hex")
	}
	sig, err := base64.StdEncoding.DecodeString(m.SignatureB64)
	if err != nil {
		return fmt.Errorf("plugin: decode signature: %w", err)
	}
	data, err := m.canonicalBytes()
	if err != nil {
		return err
	}
	if !ed25519.Verify(pub, data, sig) {
		return errors.New("plugin: signature mismatch")
	}
	return nil
}

// canonicalBytes returns YAML with the signature/pubkey fields stripped
// so the manifest signs over its own substantive content.
func (m *Manifest) canonicalBytes() ([]byte, error) {
	cp := *m
	cp.SignatureB64 = ""
	cp.PubKeyHex = ""
	return yaml.Marshal(cp)
}

func hexDecode(s string) ([]byte, error) {
	if len(s)%2 != 0 {
		return nil, errors.New("odd length")
	}
	out := make([]byte, len(s)/2)
	for i := 0; i < len(out); i++ {
		hi, err := hexNibble(s[i*2])
		if err != nil {
			return nil, err
		}
		lo, err := hexNibble(s[i*2+1])
		if err != nil {
			return nil, err
		}
		out[i] = hi<<4 | lo
	}
	return out, nil
}

func hexNibble(b byte) (byte, error) {
	switch {
	case b >= '0' && b <= '9':
		return b - '0', nil
	case b >= 'a' && b <= 'f':
		return b - 'a' + 10, nil
	case b >= 'A' && b <= 'F':
		return b - 'A' + 10, nil
	}
	return 0, fmt.Errorf("invalid hex byte %q", b)
}
