package integrity

import (
	"strings"
	"testing"
)

func TestHashReaderDeterministic(t *testing.T) {
	a, na, err := HashReader(strings.NewReader("joblantern-2026-archive"))
	if err != nil {
		t.Fatal(err)
	}
	b, nb, err := HashReader(strings.NewReader("joblantern-2026-archive"))
	if err != nil {
		t.Fatal(err)
	}
	if a != b || na != nb {
		t.Fatalf("hash not deterministic: %s != %s", a, b)
	}
}

func TestHashReproducibilityAcrossRuns(t *testing.T) {
	// Surface a stable expectation so any change to format goes
	// through an explicit decision.
	got, _, _ := HashReader(strings.NewReader(""))
	const want = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if got != want {
		t.Fatalf("empty hash drifted: got %s want %s", got, want)
	}
}
