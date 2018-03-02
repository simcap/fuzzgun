package fuzzy

import (
	"unicode"
	"bytes"
	"fmt"
	"math/rand"
	"time"
	"strings"
	"math"
)

func Fuzz(s string) (out chan string) {
	out = make(chan string)
	tokens := tokenize(s)
	go func() {
		var b bytes.Buffer
		rand.Seed(time.Now().UnixNano())
		for {
			for _, tok := range tokens {
				b.WriteString(fuzz(tok))
			}
			out <- b.String()
			b.Reset()
		}
	}()

	return out
}

func fuzz(tok string) string {
	r := rune(tok[0])
	if unicode.IsLetter(rune(r)){
		return fuzzAlpha(tok)
	} else if isSep(r){
		return fuzzSep(tok)
	} else if unicode.IsDigit(r) {
		return fuzzDigit(tok)
	}

	return tok
}

func fuzzAlpha(s string) string {
	switch rand.Intn(5) {
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

func fuzzSep(s string) string {
	return s
}

func fuzzDigit(s string) string {
	switch rand.Intn(7) {
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
	case 6:
		return fmt.Sprintf("%d", math.MaxInt64)
	case 7:
		return fmt.Sprintf("%d", math.MinInt64)
	}
	return s

}

func isSep(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}


func tokenize(s string) (tokens []string){
	if len(s) < 1 {
		return
	}

	var b bytes.Buffer
	var last rune
	for i, r := range s {
		if unicode.IsLetter(r) {
			if	unicode.IsLetter(last) || i == 0{
				last = r
				b.WriteRune(r)
				continue
			} else {
				tokens = append(tokens, b.String())
				b.Reset()
				b.WriteRune(r)
			}
		}
		if unicode.IsDigit(r) {
			if	unicode.IsDigit(last) || i == 0{
				last = r
				b.WriteRune(r)
				continue
			} else {
				tokens = append(tokens, b.String())
				b.Reset()
				b.WriteRune(r)
			}
		}
		if isSep(r) {
			if	isSep(last) || i == 0{
				last = r
				b.WriteRune(r)
				continue
			} else {
				tokens = append(tokens, b.String())
				b.Reset()
				b.WriteRune(r)
			}
		}
		last = r
	}

	if b.Len() > 0 {
		tokens = append(tokens, b.String())
	}

	return
}

