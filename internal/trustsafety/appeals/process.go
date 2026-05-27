// Package appeals is the council's entry point for entity appeals.
// An appeal flows through stages: intake, evidence-gathering, vote.
package appeals

import (
	"errors"
	"strings"
	"time"
)

// Stage is the appeal stage.
type Stage string

const (
	StageIntake          Stage = "intake"
	StageEvidence        Stage = "evidence"
	StageVote            Stage = "vote"
	StageDecided         Stage = "decided"
)

// Appeal is the council appeal record.
type Appeal struct {
	ID           string
	Subject      string
	FiledBy      string
	FiledAt      time.Time
	Stage        Stage
	StageAt      time.Time
	Decision     string
	DecidedAt    time.Time
	Recused      []string
	Votes        map[string]string // seatID -> "uphold"|"reject"|"abstain"
}

// Promote advances the appeal to the next stage. Returns an error if
// the appeal is already decided.
func (a *Appeal) Promote() error {
	if a.Stage == StageDecided {
		return errors.New("appeal already decided")
	}
	a.Stage = nextStage(a.Stage)
	a.StageAt = time.Now().UTC()
	return nil
}

func nextStage(s Stage) Stage {
	switch s {
	case StageIntake:
		return StageEvidence
	case StageEvidence:
		return StageVote
	case StageVote:
		return StageDecided
	}
	return StageDecided
}

// Validate checks the minimum content for intake.
func (a Appeal) Validate() error {
	if strings.TrimSpace(a.Subject) == "" {
		return errors.New("subject required")
	}
	if strings.TrimSpace(a.FiledBy) == "" {
		return errors.New("filedBy required")
	}
	return nil
}
