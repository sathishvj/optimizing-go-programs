package main

import "fmt"

func main() {
	a := [5]int{1, 2, 3, 4, 5}
	s := a[0:3]
	fmt.Println(a, s)

	s = append(s, 9)
	fmt.Println(a, s)

	s = append(s, 19)
	fmt.Println(a, s)

	s = append(s, 99)
	fmt.Println(a, s)
}
