package main

import "testing"

const size = 2000000

func f() ([size]int, [size]int) {
	a := [size]int{}
	b := [size]int{}
	a[19] = 100
	return a, b
}

func f2() [size]int {
	a := [size]int{}
	a[19] = 100
	return a
}

func BenchmarkHelloWorld(b *testing.B) {
	// t.Fatal("not implemented")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a := f2()
		_ = a
	}
}
