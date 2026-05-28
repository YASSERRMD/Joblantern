// Package removal documents the for-cause removal mechanism for
// council members.
package removal

// Cause is the named removal reason.
type Cause string

const (
	CauseBreachCOI       Cause = "undeclared-conflict"
	CauseAbsenteeism     Cause = "persistent-absenteeism"
	CauseConfidentiality Cause = "confidentiality-breach"
	CauseMisconduct      Cause = "documented-misconduct"
)

// Threshold returns the council vote share required to remove a seat.
func Threshold(c Cause) float64 {
	switch c {
	case CauseAbsenteeism:
		return 0.5
	case CauseBreachCOI, CauseConfidentiality, CauseMisconduct:
		return 0.6667
	}
	return 0.75
}
