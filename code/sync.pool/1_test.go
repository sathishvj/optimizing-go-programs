package main

import (
	"bytes"
	"testing"
)

func Benchmark_f1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f1()
	}
}

func f1() {
	s := &bytes.Buffer{}
	s.Write([]byte("dirty"))

	return
}
