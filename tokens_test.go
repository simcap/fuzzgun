package fuzzgun

import (
	"net"
	"net/mail"
	"reflect"
	"strconv"
	"testing"
)

func tok(s string) *token {
	return &token{s: s}
}

func allTok(s ...string) (toks []*token) {
	for _, a := range s {
		toks = append(toks, tok(a))
	}
	return
}

func TestGroupToken(t *testing.T) {
	got := groupByShifting([]*token{tok("1"), tok("2"), tok("3"), tok("4")})
	want := [][]*token{
		allTok("1"), allTok("2"), allTok("3"), allTok("4"),
		allTok("1", "2"), allTok("2", "3"), allTok("3", "4"),
		allTok("1", "2", "3"), allTok("2", "3", "4"),
		allTok("1", "2", "3", "4"),
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got\n%#v\nwant\n%#v", got, want)
	}
}

func TestDetectTypes(t *testing.T) {
	t.Run("detecting IP", func(t *testing.T) {
		s := "http://127.0.0.1"
		groups := groupByShifting(tokenize(s))

		for _, g := range groups {
			s := join(g)
			if ip := net.ParseIP(s); ip != nil {
				return
			}
		}
		t.Fatal("should have detected IP")
	})

	t.Run("detecting email", func(t *testing.T) {
		s := "http://bob@mail.com/users"
		groups := groupByShifting(tokenize(s))

		for _, g := range groups {
			s := join(g)
			if _, err := mail.ParseAddress(s); err == nil {
				return
			}
		}
		t.Fatal("should have detected email")
	})

	t.Run("detecting unix timestamp", func(t *testing.T) {
		s := "aunixtimestamp1520698065  withnano1520698065172897338"
		var found int
		for _, t := range tokenize(s) {
			if t.typ == numTok {
				if len(t.s) < 11 {
					if _, err := strconv.ParseInt(t.s, 10, 64); err == nil {
						found++
					}
				} else {
					if _, err := strconv.ParseInt(t.s, 10, 64); err == nil {
						found++
					}
				}
			}
		}
		if found != 2 {
			t.Fatal("should have detected 2 timestamps")
		}
	})
}

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
		{in: "any12345 or56789", out: []*token{newATok("any", 0), newNTok("12345", 1), newSTok(" ", 2), newATok("or", 3), newNTok("56789", 4)}},
	}

	for i, tcase := range tcases {
		actual := tokenize(tcase.in)
		if got, want := actual, tcase.out; !reflect.DeepEqual(got, want) {
			t.Fatalf("case %d\ngot  %v\n\nwant %v\n", i+1, got, want)
		}
	}

}
