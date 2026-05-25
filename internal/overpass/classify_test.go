package overpass_test

import (
	"testing"

	"github.com/yasserrmd/joblantern/internal/overpass"
)

func TestClassify(t *testing.T) {
	cases := []struct {
		name string
		els  []overpass.Element
		want string
	}{
		{"empty", nil, "unknown"},
		{"residential", []overpass.Element{{Tags: map[string]string{"building": "apartments"}}}, "residential"},
		{"commercial", []overpass.Element{{Tags: map[string]string{"office": "company"}}}, "commercial"},
		{"mixed", []overpass.Element{
			{Tags: map[string]string{"building": "apartments"}},
			{Tags: map[string]string{"shop": "convenience"}},
		}, "mixed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := &overpass.Response{Elements: tc.els}
			got := overpass.Classify(r)
			if got.PrimaryType != tc.want {
				t.Errorf("primary=%q want %q", got.PrimaryType, tc.want)
			}
		})
	}
}
