package main

import "testing"

func Benchmark_unmarshaljson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		unmarshaljsonFn()
	}
}

func Benchmark_easyjson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		easyjsonFn()
	}
}
