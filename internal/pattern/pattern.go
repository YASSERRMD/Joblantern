// Package pattern is a deterministic, rule-based red-flag classifier
// for recruitment-listing text. Rules are loaded from a YAML pack so
// they can be reviewed and evolved without recompiling.
package pattern

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

//go:embed rules.yaml
var defaultRulesYAML []byte

// RuleSpec is one entry in rules.yaml.
type RuleSpec struct {
	Code        string   `yaml:"code"`
	Description string   `yaml:"description"`
	Weight      float64  `yaml:"weight"`
	Patterns    []string `yaml:"patterns"`
}

// CompiledRule holds the compiled regex per pattern.
type CompiledRule struct {
	RuleSpec
	Regexps []*regexp.Regexp
}

// RulePack is the parsed + compiled rule set.
type RulePack struct {
	Rules []CompiledRule
}

// Hit is a single matched red flag.
type Hit struct {
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Weight      float64 `json:"weight"`
	Span        string  `json:"span"`
}

// Result is the output of Analyse.
type Result struct {
	RedFlags       []Hit   `json:"red_flags"`
	CompositeScore float64 `json:"composite_score"`
}

// DefaultPack loads the embedded rules.yaml.
func DefaultPack() (*RulePack, error) {
	return LoadPack(defaultRulesYAML)
}

// LoadPack parses + compiles a YAML rule pack.
func LoadPack(data []byte) (*RulePack, error) {
	var raw struct {
		Rules []RuleSpec `yaml:"rules"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("rule pack yaml: %w", err)
	}
	rp := &RulePack{Rules: make([]CompiledRule, 0, len(raw.Rules))}
	for _, r := range raw.Rules {
		if r.Code == "" {
			return nil, fmt.Errorf("rule missing code")
		}
		if r.Weight <= 0 || r.Weight > 1 {
			return nil, fmt.Errorf("rule %q weight out of [0,1]: %v", r.Code, r.Weight)
		}
		cr := CompiledRule{RuleSpec: r, Regexps: make([]*regexp.Regexp, 0, len(r.Patterns))}
		for _, p := range r.Patterns {
			re, err := regexp.Compile("(?i)" + p)
			if err != nil {
				return nil, fmt.Errorf("rule %q regex %q: %w", r.Code, p, err)
			}
			cr.Regexps = append(cr.Regexps, re)
		}
		rp.Rules = append(rp.Rules, cr)
	}
	return rp, nil
}

// Analyse runs every rule's regex against text and aggregates hits.
// CompositeScore is the maximum weight among triggered rules (i.e. one
// strong signal dominates rather than many weak signals stacking).
func (rp *RulePack) Analyse(text string) Result {
	var out Result
	if text == "" {
		return out
	}
	for _, r := range rp.Rules {
		for _, re := range r.Regexps {
			loc := re.FindStringIndex(text)
			if loc == nil {
				continue
			}
			span := text[loc[0]:loc[1]]
			out.RedFlags = append(out.RedFlags, Hit{
				Code: r.Code, Description: r.Description, Weight: r.Weight, Span: span,
			})
			if r.Weight > out.CompositeScore {
				out.CompositeScore = r.Weight
			}
			break // one match per rule is enough
		}
	}
	return out
}

// LanguageMismatchCheck returns true when the text contains scripts
// that conflict with the claimed jurisdiction.
//
// In v1 we implement a tiny heuristic: jurisdictions in Latin-script
// regions (UAE, Saudi, GCC) that contain Cyrillic OR CJK characters
// are flagged. Refinements live in future rule packs.
func LanguageMismatchCheck(text, jurisdiction string) (bool, string) {
	if text == "" {
		return false, ""
	}
	latinExpected := map[string]bool{
		"AE": true, "SA": true, "BH": true, "KW": true, "OM": true, "QA": true,
		"PH": true, "IN": true, "BD": true, "PK": true, "NP": true, "LK": true,
		"US": true, "GB": true, "EU": true,
	}
	if !latinExpected[strings.ToUpper(jurisdiction)] {
		return false, ""
	}
	for _, r := range text {
		switch {
		case unicode.Is(unicode.Cyrillic, r):
			return true, "cyrillic"
		case unicode.Is(unicode.Han, r), unicode.Is(unicode.Hangul, r), unicode.Is(unicode.Hiragana, r), unicode.Is(unicode.Katakana, r):
			return true, "cjk"
		}
	}
	return false, ""
}
