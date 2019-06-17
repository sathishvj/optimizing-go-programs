package mydefer

import (
	"sync"
	"testing"
)

type T struct {
	mu sync.Mutex
	n  int64
}

var t T

func (t *T) CounterA() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.n
}

func (t *T) CounterB() (count int64) {
	t.mu.Lock()
	count = t.n
	t.mu.Unlock()
	return
}

func (t *T) IncreaseA() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.n++
}

func (t *T) IncreaseB() {
	t.mu.Lock()
	t.n++ // this line will not panic for sure
	t.mu.Unlock()
}

func Benchmark_CounterA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t.CounterA()
	}
}

func Benchmark_CounterB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t.CounterB()
	}
}

func Benchmark_IncreaseA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t.IncreaseA()
	}
}

func Benchmark_IncreaseB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t.IncreaseB()
	}
}
