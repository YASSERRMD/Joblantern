// Package learning builds anonymised training datasets from accumulated
// verdicts + feedback, and computes how well each pattern rule
// correlates with confirmed-scam outcomes.
//
// Joblantern does **not** retrain a foundation model. The learning
// loop is deliberately conservative: it produces (a) labelled examples
// for downstream researchers to use under their own ethical review and
// (b) per-rule effectiveness numbers a human reviewer uses to adjust
// rule weights in `internal/pattern/rules.yaml`.
package learning

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// Verdict + feedback joined into one labelled record.
type Labelled struct {
	VerificationID string
	CompletedAt    time.Time
	Country        string
	Risk           string   // engine verdict: green/yellow/red
	PatternCodes   []string // rules that fired during verification
	Outcome        string   // user feedback: confirmed_scam | confirmed_legit | unsure | empty
}

// Anonymise scrubs anything that could identify a user. We keep:
//   - verification id (uuid; not joined to a user account here)
//   - UTC date only (no time-of-day to reduce timing-fingerprint risk)
//   - country, risk, pattern codes, feedback outcome
type AnonymisedRow struct {
	VerificationID string
	Date           string // YYYY-MM-DD
	Country        string
	Risk           string
	PatternCodes   string // comma-joined
	Outcome        string
}

// ExportCSV writes the anonymised dataset to w.
func ExportCSV(w io.Writer, rows []Labelled) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()
	if err := cw.Write([]string{"verification_id", "date", "country", "risk", "pattern_codes", "outcome"}); err != nil {
		return err
	}
	for _, r := range rows {
		a := anonymise(r)
		if err := cw.Write([]string{
			a.VerificationID, a.Date, a.Country, a.Risk, a.PatternCodes, a.Outcome,
		}); err != nil {
			return err
		}
	}
	return nil
}

func anonymise(r Labelled) AnonymisedRow {
	codes := append([]string(nil), r.PatternCodes...)
	sort.Strings(codes)
	return AnonymisedRow{
		VerificationID: r.VerificationID,
		Date:           r.CompletedAt.UTC().Format("2006-01-02"),
		Country:        r.Country,
		Risk:           r.Risk,
		PatternCodes:   strings.Join(codes, ","),
		Outcome:        r.Outcome,
	}
}

// RuleEffectiveness measures how well each pattern rule correlates
// with `confirmed_scam` feedback. Returns a map from rule code to
// EffectivenessScore.
type EffectivenessScore struct {
	Code            string  `json:"code"`
	NTotal          int     `json:"n_total"`          // how many times the rule fired
	NConfirmedScam  int     `json:"n_confirmed_scam"` // of those, how many got confirmed_scam feedback
	NConfirmedLegit int     `json:"n_confirmed_legit"`
	Precision       float64 `json:"precision"`      // confirmed_scam / (confirmed_scam + confirmed_legit) when both observed
	Recommendation  string  `json:"recommendation"` // "keep", "review", "consider_removal"
}

// MinEvidence is the minimum sample size for a recommendation other
// than "keep".
const MinEvidence = 30

// Effectiveness aggregates per-rule outcomes.
func Effectiveness(rows []Labelled) []EffectivenessScore {
	stats := map[string]*EffectivenessScore{}
	for _, r := range rows {
		for _, code := range r.PatternCodes {
			s := stats[code]
			if s == nil {
				s = &EffectivenessScore{Code: code}
				stats[code] = s
			}
			s.NTotal++
			switch r.Outcome {
			case "confirmed_scam":
				s.NConfirmedScam++
			case "confirmed_legit":
				s.NConfirmedLegit++
			}
		}
	}
	out := make([]EffectivenessScore, 0, len(stats))
	for _, s := range stats {
		denom := s.NConfirmedScam + s.NConfirmedLegit
		if denom > 0 {
			s.Precision = float64(s.NConfirmedScam) / float64(denom)
		}
		s.Recommendation = recommend(*s)
		out = append(out, *s)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Precision != out[j].Precision {
			return out[i].Precision > out[j].Precision
		}
		return out[i].NTotal > out[j].NTotal
	})
	return out
}

func recommend(s EffectivenessScore) string {
	if s.NTotal < MinEvidence {
		return "keep"
	}
	switch {
	case s.Precision >= 0.7:
		return "keep"
	case s.Precision >= 0.4:
		return "review"
	default:
		return "consider_removal"
	}
}

// SummaryText returns a short human-readable report.
func SummaryText(scores []EffectivenessScore) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Joblantern rule effectiveness report\n\n")
	fmt.Fprintf(&b, "%-32s  %6s  %6s  %6s  %6s  %s\n", "code", "n", "scam", "legit", "prec", "rec")
	for _, s := range scores {
		fmt.Fprintf(&b, "%-32s  %6d  %6d  %6d  %6.2f  %s\n",
			s.Code, s.NTotal, s.NConfirmedScam, s.NConfirmedLegit, s.Precision, s.Recommendation)
	}
	return b.String()
}
