package main

type T struct {
	a int
}

func h() *T {
	return &T{}
}

func main() {
	h()
}
