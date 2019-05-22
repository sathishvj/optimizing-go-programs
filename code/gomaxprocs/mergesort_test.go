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