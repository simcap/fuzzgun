package fuzzgun

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

func Fuzz(ctx context.Context, s string, maxIter ...int) <-chan string {
	rand.Seed(time.Now().UnixNano())

	var max int
	if len(maxIter) > 0 {
		max = maxIter[0]
	}

	out := make(chan string)
	toks := tokenize(s)
	groups := group(toks, 4)
	unique := map[string]struct{}{}

	var total int
	go func() {
		for {
			select {
			default:
				for _, tuple := range groups {
					result := make([]string, len(toks))
					src := stringArr(toks)
					copy(result, src)
					var c []int
					for _, t := range tuple {
						r := fuzz(t)
						result[t.pos] = r
						c = append(c, t.pos)
					}
					fuzzed := strings.Join(result, "")
					if _, ok := unique[fuzzed]; ok {
						continue
					} else {
						unique[fuzzed] = struct{}{}
						out <- fuzzed
						total++
						if max > 0 && total == max {
							close(out)
							return
						}
					}
				}
			case <-ctx.Done():
				close(out)
				return
			}
		}
	}()

	return out
}

func fuzz(tok *token) string {
	switch tok.typ {
	case alpha:
		return mutateAlpha(tok.s)
	case digit:
		return mutateDigit(tok.s)
	case separator:
		return mutateSep(tok.s)
	}
	return tok.s
}

type tokenType int

const (
	alpha tokenType = iota
	digit
	separator
)

type token struct {
	pos int
	s   string
	typ tokenType
}

func stringArr(t []*token) (out []string) {
	for _, t := range t {
		out = append(out, t.s)
	}
	return
}

func group(tokens []*token, groupFactor ...int) (groups [][]*token) {
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

func mutateAlpha(s string, level ...int) string {
	var m int
	if len(level) > 0 {
		m = level[0]
	} else {
		m = rand.Intn(6)
	}

	switch m {
	case 0:
		return s
	case 1:
		return strings.Repeat(s, 10)
	case 2:
		return "\x00"
	case 3:
		return "\x00" + s + "\x00"
	case 4:
		return "\x00" + s + "\x00"
	case 5:
		return ""
	}
	return s
}

func mutateDigit(s string) string {
	switch rand.Intn(2) {
	case 0:
		switch rand.Intn(5) {
		case 0:
			return s
		case 1:
			return ""
		case 2:
			return "-" + s
		case 3:
			return strings.Repeat("9", len(s)+1)
		case 4:
			return fmt.Sprintf("%s.%s", s, s)
		}
	case 1:
		switch rand.Intn(2) {
		case 0:
			return fmt.Sprintf("%d", math.MaxInt64)
		case 1:
			return fmt.Sprintf("%d", math.MinInt64)
		}
	case 2:
		return "\x00"
	case 3:
	}
	return s
}

var chars = []byte{33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 58, 59, 60, 61, 62, 63, 64, 91, 92, 93, 94, 95, 96, 123, 124, 125, 126}

func mutateSep(s string) string {
	switch rand.Intn(2) {
	case 0:
		switch rand.Intn(3) {
		case 0:
			return s
		case 1:
			return ""
		case 2:
			return s + s
		}
	case 1:
		switch rand.Intn(1) {
		case 0:
			return fmt.Sprintf("%c", chars[rand.Intn(len(chars))])
		}
	}
	return s
}

func isSep(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}

func tokenize(s string) (tokens []*token) {
	if len(s) < 1 {
		return
	}

	var b bytes.Buffer
	var last rune
	for i, r := range s {
		if unicode.IsLetter(r) {
			if unicode.IsLetter(last) || i == 0 {
				last = r
				b.WriteRune(r)
				continue
			} else {
				tokens = append(tokens, &token{s: b.String(), typ: alpha})
				b.Reset()
				b.WriteRune(r)
			}
		}
		if unicode.IsDigit(r) {
			if unicode.IsDigit(last) || i == 0 {
				last = r
				b.WriteRune(r)
				continue
			} else {
				tokens = append(tokens, &token{s: b.String(), typ: digit})
				b.Reset()
				b.WriteRune(r)
			}
		}
		if isSep(r) {
			if isSep(last) || i == 0 {
				last = r
				b.WriteRune(r)
				continue
			} else {
				tokens = append(tokens, &token{s: b.String(), typ: separator})
				b.Reset()
				b.WriteRune(r)
			}
		}
		last = r
	}

	if b.Len() > 0 {
		a := b.String()
		tokens = append(tokens, &token{s: a, typ: detectTyp(a)})
	}

	for i, tok := range tokens {
		tok.pos = i
	}

	return
}

func detectTyp(s string) tokenType {
	r := rune(s[0])
	for _, a := range s {
		if unicode.IsLetter(a) {
			return alpha
		} else if isSep(r) {
			return separator
		} else if unicode.IsDigit(r) {
			return digit
		}
	}
	return separator
}
