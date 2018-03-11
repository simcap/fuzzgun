package fuzzgun

import (
	"fmt"
	"strings"
)

type sequitur struct {
	allRules   []string
	finalRules []string
	idx        int
}

type symbol struct {
	s   string
	typ int
}

func (s symbol) String() string {
	if s.typ > 0 {
		return fmt.Sprintf("nonterm[%s]", s.s)
	}
	return fmt.Sprintf("term[%s]", s.s)
}

type digram [2]symbol

func stringToSymbols(s string) (out []symbol) {
	for _, a := range s {
		out = append(out, symbol{s: string(a), typ: 1})
	}
	return
}

func (seq *sequitur) parse(symbols []symbol) []symbol {
	digrams := make(map[digram]int)
	for i := 0; i <= len(symbols)-2; i++ {
		var d digram
		copy(d[:], symbols[i:i+2])
		digrams[d]++
	}

	rules := make(map[digram]struct{})
	for digram, count := range digrams {
		if count > 1 {
			rules[digram] = struct{}{}
		}
	}

	if len(rules) < 1 {
		return compress(symbols)
	}

	// Create new symbols replacing with rules
	var newSymbols []symbol
	for i := 0; i < len(symbols); {
		if i < len(symbols)-1 {
			slice := symbols[i : i+2]
			var d digram
			copy(d[:], slice)
			if _, ok := rules[d]; ok {
				var s []string
				s = append(s, d[0].s, d[1].s)
				term := symbol{typ: 0, s: strings.Join(s, "")}
				newSymbols = append(newSymbols, term)
				i = i + 2
			} else {
				newSymbols = append(newSymbols, symbols[i])
				i++
			}
		} else {
			newSymbols = append(newSymbols, symbols[i])
			i++
		}
	}

	return seq.parse(compress(newSymbols))
}

func compress(arr []symbol) []symbol {
	out := make([]symbol, 0)
	for _, s := range arr {
		if len(out) > 0 && s == out[len(out)-1] {
			continue
		} else {
			out = append(out, s)
		}
	}
	return out
}
