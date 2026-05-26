// Package transparency builds anonymised aggregate views of verdicts
// for a public dashboard and a CSV/JSONL export feed. The intent is
// for journalists, regulators, and policy researchers to see the
// shape of recruitment fraud without seeing individuals.
//
// Privacy posture:
//   - No row in the published dataset references an individual user.
//   - Small cells (count < `minCell`) are dropped entirely or fuzzed
//     with discrete Laplace noise so single-respondent re-identification
//     is impractical.
//   - Aggregates are computed nightly off the warm `verifications`
//     table and stored in a materialised snapshot. The public endpoints
//     read the snapshot, never live data.
package transparency

import (
	"math/rand/v2"
	"sort"
	"time"
)

// Row is one (country, risk) bucket for a time window.
type Row struct {
	Date    string `json:"date"`    // YYYY-MM-DD UTC
	Country string `json:"country"` // ISO 3166-1 alpha-2 or "ZZ"
	Risk    string `json:"risk"`    // green | yellow | red
	Count   int    `json:"count"`
	Fuzzed  bool   `json:"fuzzed"` // true if noise was applied
}

// Verdict is the minimal input record.
type Verdict struct {
	CompletedAt time.Time
	Country     string
	Risk        string
}

// Aggregator is the deterministic-when-no-noise rollup.
type Aggregator struct {
	MinCell    int     // counts strictly below this are dropped or fuzzed
	NoiseScale float64 // discrete-Laplace b parameter; 0 = no noise
	NoiseRNG   *rand.Rand
}

// New returns an aggregator with sensible defaults: drop cells < 5,
// add b=1 discrete-Laplace noise to all kept cells, so small-cohort
// re-identification is bounded.
func New() *Aggregator {
	return &Aggregator{
		MinCell:    5,
		NoiseScale: 1.0,
		NoiseRNG:   rand.New(rand.NewPCG(1, 2)),
	}
}

// Aggregate groups verdicts by (date, country, risk) and applies the
// noise / suppression policy.
func (a *Aggregator) Aggregate(in []Verdict) []Row {
	buckets := map[string]int{}
	for _, v := range in {
		if v.CompletedAt.IsZero() {
			continue
		}
		d := v.CompletedAt.UTC().Format("2006-01-02")
		country := v.Country
		if country == "" {
			country = "ZZ"
		}
		risk := v.Risk
		if risk == "" {
			continue
		}
		key := d + "|" + country + "|" + risk
		buckets[key]++
	}

	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]Row, 0, len(keys))
	for _, k := range keys {
		date, country, risk := split3(k, '|')
		count := buckets[k]
		if count < a.MinCell {
			// Drop cells below the minimum to remove single-respondent risk.
			continue
		}
		row := Row{Date: date, Country: country, Risk: risk, Count: count}
		if a.NoiseScale > 0 && a.NoiseRNG != nil {
			row.Count = max0(count + a.laplaceNoise())
			row.Fuzzed = true
		}
		out = append(out, row)
	}
	return out
}

// laplaceNoise returns an integer drawn from a discrete Laplace
// distribution centred on zero. We synthesise it from two
// independent geometric samples (Geometric(1-e^{-1/b})). Sufficient
// for the small noise scales (b ≤ 2) we target.
func (a *Aggregator) laplaceNoise() int {
	if a.NoiseScale <= 0 {
		return 0
	}
	p := 1.0 - exp(-1.0/a.NoiseScale)
	g1 := geometric(a.NoiseRNG, p)
	g2 := geometric(a.NoiseRNG, p)
	return g1 - g2
}

func geometric(r *rand.Rand, p float64) int {
	if p <= 0 || p >= 1 {
		return 0
	}
	// Inverse-CDF sampling: floor(log(U)/log(1-p)).
	u := r.Float64()
	if u <= 0 {
		return 0
	}
	return int(logFloor(u) / logFloor(1.0-p))
}

// Small helpers that avoid math imports to keep the file dep-free.
func exp(x float64) float64 {
	// 4-term Taylor; we only feed it small negative numbers (|x| < 1).
	return 1 + x + x*x/2 + x*x*x/6 + x*x*x*x/24
}

func logFloor(x float64) float64 {
	// Crude natural log via series for x in (0, 2). Adequate for noise
	// generation at b ≈ 1-2.
	if x <= 0 {
		return -1e9
	}
	y := (x - 1) / (x + 1)
	y2 := y * y
	return 2 * (y + y2*y/3 + y2*y2*y/5)
}

func max0(n int) int {
	if n < 0 {
		return 0
	}
	return n
}

func split3(s string, sep byte) (string, string, string) {
	first := -1
	second := -1
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			if first < 0 {
				first = i
			} else if second < 0 {
				second = i
				break
			}
		}
	}
	if first < 0 || second < 0 {
		return s, "", ""
	}
	return s[:first], s[first+1 : second], s[second+1:]
}
