package fuzzgun

import (
	"reflect"
	"testing"
)

func TestSequitur(t *testing.T) {
	tcases := []struct {
		in  string
		out []string
	}{
		{in: "a", out: nil}, {in: "ab", out: nil}, {in: "abc", out: nil},
		{in: "abcbc", out: []string{"bc"}},
		{in: "abaaba", out: []string{"aba"}},
		{in: "<abc><abc>", out: []string{"<abc>"}},
	}

	for i, tcase := range tcases {
		seq := &sequitur{}
		seq.run(tcase.in)
		if got, want := seq.finalRules, tcase.out; !reflect.DeepEqual(got, want) {
			t.Errorf("case %d\ngot %v\n\nwant %v\n(all rules: %s)", i+1, got, want, seq.allRules)
		}
	}
}
