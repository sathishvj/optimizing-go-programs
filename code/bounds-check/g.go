package main

import "fmt"

func g1(b []byte, v uint32) {
	b[0] = byte(v + 48)
	b[1] = byte(v + 49)
	b[2] = byte(v + 50)
	b[3] = byte(v + 51)
	fmt.Println(b)
}

func g2(b []byte, v uint32) {
	b[3] = byte(v + 51)
	b[0] = byte(v + 48)
	b[1] = byte(v + 49)
	b[2] = byte(v + 50)
	fmt.Println(b)
}

func main() {
	b := make([]byte, 4)
	g1(b, 10)
	g2(b, 10)
}
