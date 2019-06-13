package main

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var s []string

func RandString_Sequential() {
	for i := 0; i < 1000; i++ {
		s = append(s, RandString(100))
	}
}

func Benchmark_Sequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandString_Sequential()
	}
}

func RandString_Concurrent() {
	for i := 0; i < 100000; i++ {
		go func() {
			s = append(s, RandString(100))
		}()
	}
}

func Benchmark_Concurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandString_Concurrent()
	}
}

var mu sync.Mutex

func RandString_Locked_Mutex() {
	for i := 0; i < 100000; i++ {
		go func() {
			mu.Lock()
			defer mu.Unlock()

			s = append(s, RandString(100))
		}()
	}
}

func Benchmark_Locked_Mutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandString_Locked_Mutex()
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	//time.Sleep(10 * time.Microsecond)
	return string(b)
}
