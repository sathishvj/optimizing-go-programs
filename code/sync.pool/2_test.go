package main

import (
	"bytes"
	"sync"
	"testing"
)

var pool2 = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func Benchmark_f2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f2()
	}
}

func f2() {
	// When getting from a Pool, you need to cast
	s := pool2.Get().(*bytes.Buffer)
	// We write to the object
	s.Write([]byte("dirty"))
	// Then put it back
	pool2.Put(s)

	return
}
