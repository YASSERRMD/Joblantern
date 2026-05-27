// Package combined joins an edu verdict with the post-graduation job
// offer commonly attached to the same scam ("study + work, guaranteed
// placement"). Either side red-flagging the other escalates the
// overall band.
package combined

// View is the merged output for the user.
type View struct {
	EduBand     string
	JobBand     string
	OverallBand string
	Notes       []string
}

// Merge returns the worst band, plus a bump if the bundle is the
// "graduate-and-work" pattern (job tied to specific employer at
// the time of admission).
func Merge(edu, job string, gradAndWork bool) View {
	worst := worstOf(edu, job)
	if gradAndWork && worst == "yellow" {
		worst = "red"
	}
	notes := []string{}
	if gradAndWork {
		notes = append(notes, "Admission and employment are bundled — a common visa-mill scam pattern.")
	}
	return View{EduBand: edu, JobBand: job, OverallBand: worst, Notes: notes}
}

func worstOf(a, b string) string {
	rank := map[string]int{"green": 0, "yellow": 1, "red": 2}
	if rank[a] >= rank[b] {
		return a
	}
	return b
}
