package main

import (
	"fmt"
	"strconv"
	"testing"
)

func fmtFn(i int) string {
	return fmt.Sprintf("%d", i)
}

func Benchmark_fmtFn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmtFn(1234)
	}
}

func strconvFn(i int) string {
	return strconv.Itoa(i)
}
func Benchmark_strconvFn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconvFn(1234)
	}
}
