package main

import (
	"encoding/json"
	"sync"
	"testing"
)

type Book2 struct {
	Author string
	Title  string
	ISBN   string
}

var bookPool = sync.Pool{
	New: func() interface{} {
		return &Book2{}
	},
}

func write2(a, t string) {
	b := bookPool.Get().(*Book2)
	b.Author = a
	b.Title = t
	b.ISBN = "abcd"
	data, _ := json.Marshal(b)
	_ = data

	bookPool.Put(b)
}

func Benchmark_write2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		write2("harry", "rowling")
	}
}
