package audit

import "testing"

func TestHashChainTamperEvident(t *testing.T) {
	l := &Log{}
	for _, e := range []Entry{
		{Regulator: "MOL-PH", Subject: "blacklist:1", Action: "push"},
		{Regulator: "MOL-PH", Subject: "stats", Action: "fetch"},
		{Regulator: "MOL-PH", Subject: "blacklist:2", Action: "push"},
	} {
		if _, err := l.Append(e); err != nil {
			t.Fatalf("append: %v", err)
		}
	}
	es := l.Entries()
	if len(es) != 3 {
		t.Fatalf("want 3 entries, got %d", len(es))
	}
	for i := 1; i < len(es); i++ {
		if es[i].PrevHash != es[i-1].Hash {
			t.Errorf("chain broken at index %d", i)
		}
	}
}

func TestMockRegulatorFlow(t *testing.T) {
	l := &Log{}
	if _, err := l.Append(Entry{Regulator: "EMB-AE", Subject: "whitelist:42", Action: "push"}); err != nil {
		t.Fatalf("mock regulator integration failed: %v", err)
	}
}
