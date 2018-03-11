package fuzzgun

import (
	"testing"
)

func TestDisplaySequitur(t *testing.T) {
	tcases := []string{
		"a", "ab", "abc",
		"abcbc", "abaaba",
		"aaabdaaabac", "<abc><abc>",
		"<abc>html<abc>",
		"the mighty the tall",
		"2.3.2.4.",
	}

	for _, s := range tcases {
		seq := &sequitur{}
		symbols := stringToSymbols(s)
		t.Log("case", s)
		t.Log("rules:", seq.parse(symbols))
		t.Log()
	}
}
