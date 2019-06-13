package main

import "fmt"

func main() {
	a := [5]int{1, 2, 3, 4, 5}
	s := a[0:3]
	s[0] = 11
	fmt.Println(a, s)

	fmt.Printf("%p %p\n", &a, &s)
	fmt.Printf("%p %p\n", &a[0], &s[0])
}
