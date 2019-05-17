package main

import (
	"os"
	"runtime/trace"
	"testing"
)

func Benchmark_h(b *testing.B) {
	var t *T

	f, err := os.Create("h.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		t = h()
	}

	trace.Stop()

	b.StopTimer()

	_ = t
}
