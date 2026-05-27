package lookup

import "testing"

func TestHashStripsFormatting(t *testing.T) {
	a := HashPhone("+1 (415) 555-1234")
	b := HashPhone("14155551234")
	if a != b {
		t.Fatalf("expected formatting-insensitive hash, got %q vs %q", a, b)
	}
}

func TestSimulatedIncomingCall(t *testing.T) {
	// A simulated scam number triggers a red overlay. Here we
	// just confirm the hash pipeline is deterministic for the
	// canonical test fixture used by the device-side suite.
	want := HashPhone("971555111000")
	if HashPhone("+971 55 511 1000") != want {
		t.Fatalf("hash mismatch on simulated call")
	}
}
