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
