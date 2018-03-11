package fuzzgun

import (
	"fmt"
	"testing"
)

func TestSequiturOnCorpus(t *testing.T) {
	tcases := [][]string{
		{
			"http://example.com",
			"https://example.com",
			"http://example.com?q=10",
			"http://example.com?q=10&v=9"},
		{
			"the tiny guy",
			"Has the bar", "the well", "the war",
		},
		{
			"<joe@mail.com>", "joe@mail.com",
			"Joe <joe@email.com>", "joe@email.com",
		},
	}
	for _, in := range tcases {
		seq := &sequitur{}
		symbols := stringToCorpus(in...)
		fmt.Println("rules:", seq.parse(symbols))
		fmt.Println()
	}
}

func TestDisplaySequitur(t *testing.T) {
	tcases := []string{
		"a", "ab", "abc",
		"abcbc", "abaaba",
		"aaabdaaabac", "<abc><abc>",
		"<abc>html<abc>",
		"the mighty the tall the ugly",
		"2.3.2.4.",
	}
	for _, s := range tcases {
		seq := &sequitur{}
		corpus := stringToCorpus(s)
		fmt.Println("rules:", seq.parse(corpus))
		fmt.Println()
	}
}
