package main

import (
	"context"
	"regexp"
	"runtime/pprof"
	"testing"
)

var ss = []string{
	`^[a-z]+\[[0-9]+\]$`,
	`foo.*`,
	`foo(.?)`,
	`foo.?`,
	`a(x*)b(y|z)c`,
}

func f(s string) {
	labels := pprof.Labels("pat", s)
	pprof.Do(context.Background(), labels, func(ctx context.Context) {
		// Do some work...
		r := regexp.MustCompile(s)
		_ = r

		//go update(ctx) // propagates labels in ctx.
	})
}

func bench_f(b *testing.B, s string) {
	for i := 0; i < b.N; i++ {
		f(s)
	}
}

func Benchmark_0f(b *testing.B) {
	bench_f(b, ss[0])
}

func Benchmark_1f(b *testing.B) {
	bench_f(b, ss[1])
}
