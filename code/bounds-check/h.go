package main

import "fmt"

func h1(b []byte, n int) {
	b[n+0] = byte(1) // Found IsInBounds
	b[n+1] = byte(2) // Found IsInBounds
	b[n+2] = byte(3) // Found IsInBounds
	b[n+3] = byte(4) // Found IsInBounds
	b[n+4] = byte(5) // Found IsInBounds
	b[n+5] = byte(6) // Found IsInBounds
	fmt.Println("in h1(): ", b)
}

func h2(b []byte, n int) {
	b = b[n : n+6] // Found IsSliceInBounds
	b[0] = byte(1)
	b[1] = byte(2)
	b[2] = byte(3)
	b[3] = byte(4)
	b[4] = byte(5)
	b[5] = byte(6)
	fmt.Println("in h2(): ", b)
}

func main() {
	b := make([]byte, 20)
	h1(b, 10)
	fmt.Println("in main: ", b)
	h2(b, 10)
	fmt.Println("in main: ", b)
}
