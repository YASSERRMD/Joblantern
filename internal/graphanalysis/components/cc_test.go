package components

import (
	"sort"
	"testing"
)

func TestSeven_Sibling_Companies_Two_Phones_One_Address(t *testing.T) {
	// Synthetic fixture: 7 fake companies share 2 phone numbers and
	// 1 address. The whole set must collapse into a single component.
	nodes := []string{"c1", "c2", "c3", "c4", "c5", "c6", "c7", "p1", "p2", "a1"}
	edges := []Edge{
		{"c1", "p1"}, {"c2", "p1"}, {"c3", "p1"},
		{"c4", "p2"}, {"c5", "p2"}, {"c6", "p2"},
		{"c7", "a1"}, {"c1", "a1"}, {"c4", "a1"},
	}
	groups := Find(nodes, edges)
	if len(groups) != 1 {
		t.Fatalf("expected 1 component, got %d", len(groups))
	}
	sort.Strings(groups[0])
	if len(groups[0]) != 10 {
		t.Fatalf("expected 10 nodes, got %d", len(groups[0]))
	}
}
