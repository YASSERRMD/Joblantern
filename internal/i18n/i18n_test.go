package i18n_test

import (
	"net/http/httptest"
	"testing"

	"github.com/yasserrmd/joblantern/internal/i18n"
)

func TestLoadAndLookup(t *testing.T) {
	c, err := i18n.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Lookup("en", "verify_button"); got != "Verify" {
		t.Errorf("en verify_button = %q", got)
	}
	if got := c.Lookup("ar", "verify_button"); got == "" || got == "verify_button" {
		t.Errorf("ar verify_button = %q", got)
	}
	// unknown key falls back to itself
	if got := c.Lookup("en", "totally_unknown_key"); got != "totally_unknown_key" {
		t.Errorf("fallback got %q", got)
	}
}

func TestResolveAcceptLanguage(t *testing.T) {
	c, _ := i18n.Load()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Language", "ar-AE;q=0.9, en;q=0.5")
	lang, dir := c.Resolve(r)
	if lang != "ar" || dir != "rtl" {
		t.Errorf("got %s %s", lang, dir)
	}
}

func TestParityAllCatalogs(t *testing.T) {
	c, _ := i18n.Load()
	// All catalogs should at least have the same keys English has.
	keys := []string{
		"tagline", "verify_button", "submit_again", "status", "risk",
		"reasons", "evidence", "confidence", "agent_working",
		"title_home", "title_result",
	}
	for _, lang := range []string{"en", "ar", "hi", "tl", "bn", "ur", "sw", "es"} {
		for _, k := range keys {
			if v := c.Lookup(lang, k); v == k {
				t.Errorf("catalog %s missing %q", lang, k)
			}
		}
	}
}
