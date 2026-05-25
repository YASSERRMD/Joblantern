// Package i18n provides a small message-catalog lookup with
// Accept-Language + cookie locale resolution. v1 ships catalogs for
// 8 languages; the API is intentionally narrow so adding more is a
// JSON-only change.
package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/text/language"
)

//go:embed catalogs/*.json
var catalogFS embed.FS

const cookieName = "joblantern_lang"

// Catalog is the parsed messages.json for one language.
type Catalog struct {
	Lang     string            `json:"-"`
	Dir      string            `json:"dir"` // ltr | rtl
	Messages map[string]string `json:"messages"`
}

// Catalogs is the registry.
type Catalogs struct {
	mu      sync.RWMutex
	byTag   map[language.Tag]*Catalog
	matcher language.Matcher
}

// Load reads every catalogs/*.json into memory and returns a Catalogs.
func Load() (*Catalogs, error) {
	c := &Catalogs{byTag: map[language.Tag]*Catalog{}}
	files, err := catalogFS.ReadDir("catalogs")
	if err != nil {
		return nil, err
	}
	tags := []language.Tag{}
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		lang := strings.TrimSuffix(f.Name(), ".json")
		tag, err := language.Parse(lang)
		if err != nil {
			return nil, fmt.Errorf("parse lang %s: %w", lang, err)
		}
		data, err := catalogFS.ReadFile("catalogs/" + f.Name())
		if err != nil {
			return nil, err
		}
		var cat Catalog
		if err := json.Unmarshal(data, &cat); err != nil {
			return nil, fmt.Errorf("parse %s: %w", f.Name(), err)
		}
		cat.Lang = lang
		c.byTag[tag] = &cat
		tags = append(tags, tag)
	}
	if len(tags) == 0 {
		return nil, fmt.Errorf("no catalogs loaded")
	}
	c.matcher = language.NewMatcher(tags)
	return c, nil
}

// Lookup returns the message text in `lang`, falling back to English.
func (c *Catalogs) Lookup(lang, key string) string {
	if cat := c.lookupCatalog(lang); cat != nil {
		if v, ok := cat.Messages[key]; ok {
			return v
		}
	}
	if cat := c.lookupCatalog("en"); cat != nil {
		if v, ok := cat.Messages[key]; ok {
			return v
		}
	}
	return key
}

// Resolve picks a language based on cookie > Accept-Language > English.
func (c *Catalogs) Resolve(r *http.Request) (string, string) {
	if ck, err := r.Cookie(cookieName); err == nil && ck.Value != "" {
		return c.resolveTag(language.Make(ck.Value))
	}
	tags, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if err == nil && len(tags) > 0 {
		matched, _, _ := c.matcher.Match(tags...)
		return c.resolveTag(matched)
	}
	return "en", "ltr"
}

func (c *Catalogs) resolveTag(t language.Tag) (string, string) {
	base, _ := t.Base()
	lang := base.String()
	dir := "ltr"
	if cat := c.lookupCatalog(lang); cat != nil {
		dir = cat.Dir
	}
	return lang, dir
}

func (c *Catalogs) lookupCatalog(lang string) *Catalog {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for tag, cat := range c.byTag {
		base, _ := tag.Base()
		if base.String() == lang {
			return cat
		}
	}
	return nil
}

// SetCookie writes a 1-year persistent locale cookie.
func SetCookie(w http.ResponseWriter, lang string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    lang,
		Path:     "/",
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
