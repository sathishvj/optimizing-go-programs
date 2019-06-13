package main

func main() {
	var a [5]int
	s := a[0:3]
	s = a[:3]
	s = a[3:]

	// negative indexing is not allowed
	// s = a[0:-2] // compile error
}
