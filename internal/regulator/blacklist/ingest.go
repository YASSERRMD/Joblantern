// Package blacklist accepts official blacklists published by verified
// regulators. Each entry is weighted heavily in the risk engine —
// effectively treated as ground truth unless an appeal is granted.
package blacklist

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"
	"time"
)

// Entry is one row in a regulator blacklist.
type Entry struct {
	RegulatorID string
	IssuedAt    time.Time
	EntityName  string
	EntityID    string
	Reason      string
	URL         string
}

// Weight is how heavily the risk engine factors a blacklist hit. The
// canonical risk score is in [0,100]; a hit pushes the score to at
// least this value.
const Weight = 95

// ParseCSV reads a regulator blacklist CSV (header: regulator_id,
// issued_at, entity_name, entity_id, reason, url).
func ParseCSV(r io.Reader, regulatorID string) ([]Entry, error) {
	cr := csv.NewReader(r)
	rows, err := cr.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, errors.New("empty blacklist")
	}
	var out []Entry
	for i, row := range rows[1:] {
		if len(row) < 6 {
			return nil, errors.New("malformed row")
		}
		ts, err := time.Parse(time.RFC3339, row[1])
		if err != nil {
			return nil, errors.New("row " + itoa(i+2) + ": bad issued_at")
		}
		out = append(out, Entry{
			RegulatorID: regulatorID,
			IssuedAt:    ts,
			EntityName:  strings.TrimSpace(row[2]),
			EntityID:    strings.TrimSpace(row[3]),
			Reason:      strings.TrimSpace(row[4]),
			URL:         strings.TrimSpace(row[5]),
		})
	}
	return out, nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
