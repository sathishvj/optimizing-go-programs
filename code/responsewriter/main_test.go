package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func withoutSetHeader(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello, stranger")
}

func Benchmark_withoutSetHeader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/", nil)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(withoutSetHeader)

		handler.ServeHTTP(rr, req)
	}

}

func withSetHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "hello, stranger")
}

func Benchmark_withSetHeader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/", nil)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(withSetHeader)

		handler.ServeHTTP(rr, req)
	}
}
