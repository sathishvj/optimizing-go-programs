package main

import "testing"

func Test_isGopher(t *testing.T) {

	tcs := []struct {
		in    string
		exp   bool
		expId string
	}{
		{
			"",
			false,
			"",
		},
		{
			"a@email.com",
			false,
			"",
		},
		{
			"a@golang.org",
			true,
			"a",
		},
	}

	for _, tc := range tcs {
		id, ok := isGopher(tc.in)
		if ok != tc.exp {
			t.Errorf("For input %s, expected: %t but got: %t", tc.in, tc.exp, ok)
		}
		if id != tc.expId {
			t.Errorf("For input %s, expected: %s but got: %s", tc.in, tc.expId, id)
		}
	}
}

func Benchmark_isGopher(b *testing.B) {

	tcs := []struct {
		in    string
		exp   bool
		expId string
	}{
		{
			"a@golang.org",
			true,
			"a",
		},
	}

	for i := 0; i < b.N; i++ {
		isGopher(tcs[0].in)
	}
}
