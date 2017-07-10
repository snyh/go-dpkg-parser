package lexer

import (
	"testing"
)

type TS struct {
	estr       string
	unexpected string
	expecting  []string
}

func TestParseError(t *testing.T) {
	errorErrorVerbose = true

	data := []TS{
		TS{
			"syntax error: unexpected $end, expecting HEAD or TOKEN",
			"$end",
			[]string{"HEAD", "TOKEN"},
		},
		TS{
			"syntax error: unexpected TOKEN, expecting HEAD",
			"TOKEN",
			[]string{"HEAD"},
		},

		TS{
			"syntax error: unexpected TOKEN",
			"TOKEN",
			[]string{},
		},
	}

	same := func(s1 []string, s2 []string) bool {
		if len(s1) != len(s2) {
			return false
		}
		for i, s := range s1 {
			if s != s2[i] {
				return false
			}
		}
		return true
	}

	for _, datum := range data {
		u, e := ParseError(datum.estr, errorToknames[:])
		if u != datum.unexpected || !same(e, datum.expecting) {
			t.Fatalf("In: %v Out: %v and %v\n", datum, u, e)
		}
	}
}
