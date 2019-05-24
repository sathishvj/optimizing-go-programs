package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	var data string
	if len(os.Args) == 2 {
		data = os.Args[1]
	}

	id, ok := isGopher(data)
	if !ok {
		id = "stranger"
	}
	fmt.Printf("hello, %s\n", id)
}

func isGopher(email string) (string, bool) {
	re := regexp.MustCompile("^([[:alpha:]]+)@golang.org$")
	match := re.FindStringSubmatch(email)
	if len(match) == 2 {
		return match[1], true
	}
	return "", false
}
