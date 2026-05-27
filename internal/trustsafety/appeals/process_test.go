package appeals

import (
	"testing"
	"time"
)

func TestFullAppealFlow(t *testing.T) {
	a := &Appeal{ID: "a-1", Subject: "Bright Stars Recruitment", FiledBy: "ngo:bd-1", FiledAt: time.Now(), Stage: StageIntake, Votes: map[string]string{}}
	if err := a.Validate(); err != nil {
		t.Fatalf("intake validate: %v", err)
	}
	for i := 0; i < 3; i++ {
		if err := a.Promote(); err != nil {
			t.Fatalf("promote: %v", err)
		}
	}
	if a.Stage != StageDecided {
		t.Fatalf("expected decided, got %s", a.Stage)
	}
	if err := a.Promote(); err == nil {
		t.Fatalf("expected error promoting after decided")
	}
}
