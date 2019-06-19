package main

func f(b []byte, n int) {
	_ = b[n-1]
	for i := 0; i < n; i++ {
		b[i] = 9
	}
}
