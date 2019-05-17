package main

type T struct {
	a int
}

func s() T {
	return T{}
}

func main() {
	s()
}
