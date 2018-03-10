package fuzzgun

import (
	"reflect"
	"testing"
)

func TestGroup(t *testing.T) {
	t.Skip()
	got := group([]*token{&token{s: "1"}, &token{s: "2"}, &token{s: "3"}, &token{s: "4"}})
	want := [][]string{
		{"1"}, {"2"}, {"3"}, {"4"},
		{"1", "2"}, {"2", "3"}, {"3", "4"},
		{"1", "2", "3"}, {"2", "3", "4"},
		{"1", "2", "3", "4"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got\n%#v\nwant\n%#v", got, want)
	}
}

func TestTokenize(t *testing.T) {
	tcases := []struct {
		in  string
		out []string
	}{
		{in: "", out: nil},
		{in: "a", out: []string{"a"}},
		{in: "ab", out: []string{"ab"}},
		{in: "11", out: []string{"11"}},
		{in: "abcde", out: []string{"abcde"}},
		{in: "////", out: []string{"////"}},
		{in: "...///..", out: []string{"...///.."}},
		{in: "01-02-03", out: []string{"01", "-", "02", "-", "03"}},
		{in: "01/02/2007", out: []string{"01", "/", "02", "/", "2007"}},
		{in: "2.4567", out: []string{"2", ".", "4567"}},
		{in: "127.0.0.1", out: []string{"127", ".", "0", ".", "0", ".", "1"}},
	}

	for i, tcase := range tcases {
		actual := tokenize(tcase.in)
		if got, want := stringArr(actual), tcase.out; !reflect.DeepEqual(got, want) {
			t.Fatalf("case %d: got %#v, want %#v", i+1, got, want)
		}
	}

}
