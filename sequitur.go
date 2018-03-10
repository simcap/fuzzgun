package fuzzgun

import (
	"fmt"
	"strings"
)

type sequitur struct {
	allRules   map[string]string
	finalRules []string
	idx        int
}

func (s *sequitur) nextID() string {
	s.idx++
	return fmt.Sprintf("%d", s.idx)
}

func (s *sequitur) run(input string) {
	if s.allRules == nil {
		s.allRules = make(map[string]string)
	}

	factor := 2
	digrams := make(map[string]int)
	for i := 0; i <= len(input)-factor; i++ {
		digrams[input[i:i+factor]]++
	}

	newRules := make(map[string]string)
	for digram, count := range digrams {
		if count > 1 {
			id := s.nextID()
			s.allRules[id] = digram
			newRules[digram] = id
		}
	}

	if len(newRules) < 1 {
		for i, f := range s.finalRules {
			new := f
			for id, a := range s.allRules {
				new = strings.Replace(new, id, a, -1)
			}
			s.finalRules[i] = new

		}
		return
	}

	s.finalRules = []string{}
	for r := range newRules {
		s.finalRules = append(s.finalRules, r)
	}

	result := input
	for digram, id := range newRules {
		result = strings.Replace(result, digram, id, -1)
	}

	s.run(result)
}
