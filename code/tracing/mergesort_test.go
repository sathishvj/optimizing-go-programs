package main

import "testing"

func Benchmark_mergesortv1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mergesortv1(s)
	}
}

func Benchmark_mergesortv2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mergesortv2(s)
	}
}

func Benchmark_mergesortv3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mergesortv3(s)
	}
}

func Test_mergesortv1(t *testing.T) {
	inp := []int{89, 123, 12, 9, 198, 1546, 108, 872, 93}
	exp := []int{9, 12, 89, 93, 108, 123, 198, 872, 1546}
	mergesortv1(inp)
	if inp[0] != exp[0] && inp[len(exp)-1] != exp[len(exp)-1] {
		t.Errorf("Test failed")
	}
}
