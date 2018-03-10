package fuzzgun

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type tokenType int

func (tt tokenType) String() string {
	switch tt {
	case alphaTok:
		return "alpha"
	case numTok:
		return "num"
	case sepTok:
		return "sep"
	}
	return "?"
}

const (
	alphaTok tokenType = iota
	numTok
	sepTok
)

type token struct {
	pos int
	s   string
	typ tokenType
}

func (t *token) String() string {
	return fmt.Sprintf("%stok_%d(%s)", t.typ, t.pos, t.s)
}

func newTok(s string, t tokenType) *token {
	return &token{s: s, typ: t}
}

func newATok(s string, pos int) *token { return &token{s: s, typ: alphaTok, pos: pos} }
func newNTok(s string, pos int) *token { return &token{s: s, typ: numTok, pos: pos} }
func newSTok(s string, pos int) *token { return &token{s: s, typ: sepTok, pos: pos} }

func tokenize(s string) (tokens []*token) {
	if len(s) < 1 {
		return
	}

	first, _ := utf8.DecodeRuneInString(s)
	currentType := detectTyp(first)
	if len(s) == 1 {
		return []*token{&token{s: s, typ: currentType, pos: 0}}
	}

	var b bytes.Buffer
	lastType := currentType
	for i, r := range s {
		currentType = detectTyp(r)
		if currentType != lastType {
			tokens = append(tokens, newTok(b.String(), lastType))
			b.Reset()
		}

		if i == len(s)-1 {
			b.WriteRune(r)
			tokens = append(tokens, newTok(b.String(), currentType))
			break
		} else {
			b.WriteRune(r)
		}

		lastType = currentType
	}

	for i, tok := range tokens {
		tok.pos = i
	}

	return
}

func isSep(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}

func detectTyp(r rune) tokenType {
	if unicode.IsLetter(r) {
		return alphaTok
	} else if isSep(r) {
		return sepTok
	} else if unicode.IsDigit(r) {
		return numTok
	}
	return sepTok
}

func join(toks []*token) string {
	var out []string
	for _, t := range toks {
		out = append(out, t.s)
	}
	return strings.Join(out, "")
}

func groupByShifting(tokens []*token, groupFactor ...int) (groups [][]*token) {
	l := len(tokens)
	factor := l
	if len(groupFactor) > 0 {
		factor = groupFactor[0]
	}
	if factor > l {
		factor = l
	}
	for i := 1; i <= factor; i++ {
		for j := 0; j+i <= l; j++ {
			groups = append(groups, tokens[j:j+i])
		}
	}
	return
}
