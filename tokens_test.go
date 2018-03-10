package fuzzgun

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tcases := []struct {
		in  string
		out []*token
	}{
		{in: "", out: nil},
		{in: "   ", out: []*token{newSTok("   ", 0)}},
		{in: "a", out: []*token{newATok("a", 0)}},
		{in: "ab", out: []*token{newATok("ab", 0)}},
		{in: "11", out: []*token{newNTok("11", 0)}},
		{in: "abcde", out: []*token{newATok("abcde", 0)}},
		{in: "////", out: []*token{newSTok("////", 0)}},
		{in: "...///..", out: []*token{newSTok("...///..", 0)}},
		{in: "01-02-03", out: []*token{newNTok("01", 0), newSTok("-", 1), newNTok("02", 2), newSTok("-", 3), newNTok("03", 4)}},
		{in: "01/02/2007", out: []*token{newNTok("01", 0), newSTok("/", 1), newNTok("02", 2), newSTok("/", 3), newNTok("2007", 4)}},
		{in: "2.4567", out: []*token{newNTok("2", 0), newSTok(".", 1), newNTok("4567", 2)}},
		{in: "127.0.0.1", out: []*token{newNTok("127", 0), newSTok(".", 1), newNTok("0", 2), newSTok(".", 3), newNTok("0", 4), newSTok(".", 5), newNTok("1", 6)}},
	}

	for i, tcase := range tcases {
		actual := tokenize(tcase.in)
		if got, want := actual, tcase.out; !reflect.DeepEqual(got, want) {
			t.Fatalf("case %d\ngot  %v\n\nwant %v\n", i+1, got, want)
		}
	}

}
