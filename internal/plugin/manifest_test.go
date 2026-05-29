package plugin_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/yasserrmd/joblantern/internal/plugin"
)

func TestLoadManifest_OK(t *testing.T) {
	src := `
name: joblantern.example
display_name: Example Plugin
version: 0.1.0
description: For tests.
author: Test
license: Apache-2.0
transport: stdio
command: mcp-example
attribution:
  text: "Powered by Example"
  url: "https://example/"
tools:
  - name: ping
    description: returns pong
`
	m, err := plugin.LoadManifest([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "joblantern.example" || m.License != "Apache-2.0" {
		t.Fatalf("got %+v", m)
	}
}

func TestLoadManifest_RejectsNonPermissive(t *testing.T) {
	src := `
name: joblantern.bad
version: 0.1.0
license: GPL-3.0
transport: stdio
command: mcp-bad
attribution:
  text: "x"
`
	_, err := plugin.LoadManifest([]byte(src))
	if err == nil || !strings.Contains(err.Error(), "permissive") {
		t.Fatalf("expected license rejection, got %v", err)
	}
}

func TestVerify_SigRoundTrip(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)

	src := `
name: joblantern.signed
version: 0.1.0
license: MIT
transport: stdio
command: mcp-signed
attribution:
  text: "Signed Co"
`
	m, _ := plugin.LoadManifest([]byte(src))

	// Sign the canonical bytes (which exclude signature/pubkey).
	cp := *m
	cp.SignatureB64 = ""
	cp.PubKeyHex = ""
	data, _ := yaml.Marshal(cp)
	sig := ed25519.Sign(priv, data)

	m.PubKeyHex = hex(pub)
	m.SignatureB64 = base64.StdEncoding.EncodeToString(sig)

	if err := m.Verify(plugin.TrustCommunity); err != nil {
		t.Fatalf("verify: %v", err)
	}
}

func TestVerify_UnsignedRequiresExternal(t *testing.T) {
	src := `
name: joblantern.unsigned
version: 0.1.0
license: MIT
transport: stdio
command: mcp-x
attribution:
  text: "X"
`
	m, _ := plugin.LoadManifest([]byte(src))
	if err := m.Verify(plugin.TrustCommunity); err == nil {
		t.Fatal("community trust must require signature")
	}
	if err := m.Verify(plugin.TrustExternal); err != nil {
		t.Fatalf("external trust should accept unsigned: %v", err)
	}
}

func hex(b []byte) string {
	const c = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, x := range b {
		out[i*2] = c[x>>4]
		out[i*2+1] = c[x&0x0f]
	}
	return string(out)
}
