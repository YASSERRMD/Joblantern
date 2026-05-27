// Package timemachine serves the "view verdicts as of date D" UI.
// The implementation pulls from the matching annual archive and the
// monthly snapshots that follow it.
package timemachine

import "time"

// Query is the user-facing input.
type Query struct {
	AsOf  time.Time
	Limit int
}

// Resolver maps an AsOf date to the storage paths that cover it.
type Resolver struct {
	AnnualArchives map[int]string // year -> archive URL
	MonthlyDeltas  map[string]string // "2026-05" -> delta URL
}

// Locate returns the archive URLs needed to reconstruct state at AsOf.
func (r Resolver) Locate(q Query) []string {
	if q.AsOf.IsZero() {
		return nil
	}
	year := q.AsOf.Year()
	out := []string{}
	if a, ok := r.AnnualArchives[year-1]; ok {
		out = append(out, a)
	}
	for m := time.January; m <= q.AsOf.Month(); m++ {
		key := monthKey(year, m)
		if d, ok := r.MonthlyDeltas[key]; ok {
			out = append(out, d)
		}
	}
	return out
}

func monthKey(year int, m time.Month) string {
	mm := int(m)
	mmStr := "0" + string('0'+byte(mm))
	if mm >= 10 {
		mmStr = string('0'+byte(mm/10)) + string('0'+byte(mm%10))
	}
	return itoa(year) + "-" + mmStr
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [12]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
