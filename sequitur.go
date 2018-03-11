package fuzzgun

import (
	"fmt"
	"reflect"
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

func stringToCorpus(corpus ...string) (out [][]symbol) {
	for _, c := range corpus {
		var symbols []symbol
		for _, a := range c {
			symbols = append(symbols, symbol{s: string(a), typ: 1})
		}
		out = append(out, symbols)
	}
	return
}

func (seq *sequitur) parse(corpus [][]symbol) (out []symbol) {
	rules := extractRules(corpus)

	if len(rules) < 1 {
		for _, symbols := range corpus {
			out = append(out, compressTerminals(symbols)...)
		}
		return out
	}

	var newCorpus [][]symbol
	for _, symbols := range corpus {
		newSymbols := replaceWithRules(symbols, rules)
		newCorpus = append(newCorpus, compressTerminals(newSymbols))
	}

	return seq.parse(newCorpus)
}

func replaceWithRules(symbols []symbol, rules map[digram]struct{}) (out []symbol) {
	for i := 0; i < len(symbols); {
		if i < len(symbols)-1 {
			slice := symbols[i : i+2]
			var d digram
			copy(d[:], slice)
			if _, ok := rules[d]; ok {
				var s []string
				s = append(s, d[0].s, d[1].s)
				term := symbol{typ: 0, s: strings.Join(s, "")}
				out = append(out, term)
				i = i + 2
			} else {
				out = append(out, symbols[i])
				i++
			}
		} else {
			out = append(out, symbols[i])
			i++
		}
	}
	return
}

func extractRules(corpus [][]symbol) map[digram]struct{} {
	digrams := make(map[digram]int)
	rules := make(map[digram]struct{})
	for _, symbols := range corpus {
		for i := 0; i <= len(symbols)-2; i++ {
			var d digram
			copy(d[:], symbols[i:i+2])
			digrams[d]++
		}

		for digram, count := range digrams {
			if count > 1 {
				rules[digram] = struct{}{}
			}
		}
	}
	return rules
}

func compressTerminals(arr []symbol) []symbol {
	out := make([]symbol, 0)
	for _, s := range arr {
		if s.typ < 1 && len(out) > 0 && reflect.DeepEqual(s, out[len(out)-1]) {
			continue
		} else {
			out = append(out, s)
		}
	}
	return out
}
