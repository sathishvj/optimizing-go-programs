package main

import (
	"strings"
	"testing"
)

func BeginsWith(s, pat string) bool {
	return strings.HasPrefix(s, pat)

}

func Test_BeginsWith(t *testing.T) {
	tc := []struct {
		s, pat string
		exp    bool
	}{
		{"GoLang", "Go", true},
		{"GoLang", "Java", false},
		{"GoLang is awesome", "awe", false},
		{"awesome is GoLang. - Yoda", "awe", true},
	}

	for _, tt := range tc {
		if BeginsWith(tt.s, tt.pat) != tt.exp {
			t.Fail()
		}
	}
}

func Benchmark_BeginsWith(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BeginsWith("GoLang", "Go")
	}
}

// forced allocations for benchmem
/*
func x() *string {
	s := "hello world there"
	return &s
}

func Benchmark_x(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a := x()
		*a += *a
		_ = a
	}
}
*/
