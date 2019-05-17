// go build -gcflags='-m' 1.go
// go build -gcflags='-m -l' 1.go to avoid inlining
// go build -gcflags='-m -l -m' 1.go for verbose comments.

package main

func f() {
	var i = 5
	i++
	_ = i
}

func f_returns() int {
	var i = 5
	i++
	return i
}

func f_returns_ptr() *int {
	var i = 5
	i++
	return &i
}

func main() {
	f()
	f_returns()
	f_returns_ptr()
}
