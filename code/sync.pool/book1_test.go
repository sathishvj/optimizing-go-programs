package main

import (
	"encoding/json"
	"testing"
)

type Book struct {
	Author string
	Title  string
	ISBN   string
}

func write1(a, t string) {
	b := &Book{}
	b.Author = a
	b.Title = t
	b.ISBN = "abcd"
	data, _ := json.Marshal(b)
	_ = data
}

func Benchmark_write1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		write1("harry", "rowling")
	}
}
