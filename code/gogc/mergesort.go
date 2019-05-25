// ref: https://hackernoon.com/parallel-merge-sort-in-go-fe14c1bc006

// GOGC=off go run mergesort.go v1 && go tool trace v1.trace
// GOGC=50 go run mergesort.go v1 && go tool trace v1.trace
// GOGC=100 go run mergesort.go v1 && go tool trace v1.trace
// GOGC=200 go run mergesort.go v1 && go tool trace v1.trace
package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/trace"
	"sync"
	"time"
)

const max = 1 << 11

func merge(s []int, middle int) {
	helper := make([]int, len(s))
	copy(helper, s)

	helperLeft := 0
	helperRight := middle
	current := 0
	high := len(s) - 1

	for helperLeft <= middle-1 && helperRight <= high {
		if helper[helperLeft] <= helper[helperRight] {
			s[current] = helper[helperLeft]
			helperLeft++
		} else {
			s[current] = helper[helperRight]
			helperRight++
		}
		current++
	}

	for helperLeft <= middle-1 {
		s[current] = helper[helperLeft]
		current++
		helperLeft++
	}
}

func mergesortv1(s []int) {
	len := len(s)

	if len > 1 {
		middle := len / 2

		var wg sync.WaitGroup
		wg.Add(2)

		// First half
		go func() {
			defer wg.Done()
			mergesortv1(s[:middle])
		}()

		// Second half
		go func() {
			defer wg.Done()
			mergesortv1(s[middle:])
		}()

		// Wait that the two goroutines are completed
		wg.Wait()
		merge(s, middle)
	}
}

/* Sequential */

func mergesort(s []int) {
	if len(s) > 1 {
		middle := len(s) / 2
		mergesort(s[:middle])
		mergesort(s[middle:])
		merge(s, middle)
	}
}

func mergesortv2(s []int) {
	len := len(s)

	if len > 1 {
		if len <= max { // Sequential
			mergesort(s)
		} else { // Parallel
			middle := len / 2

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				mergesortv2(s[:middle])
			}()

			go func() {
				defer wg.Done()
				mergesortv2(s[middle:])
			}()

			wg.Wait()
			merge(s, middle)
		}
	}
}

func mergesortv3(s []int) {
	len := len(s)

	if len > 1 {
		if len <= max { // Sequential
			mergesort(s)
		} else { // Parallel
			middle := len / 2

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()
				mergesortv3(s[:middle])
			}()

			mergesortv3(s[middle:])

			wg.Wait()
			merge(s, middle)
		}
	}
}

// Generates a slice of size, size filled with random numbers
func generateSlice(size int) []int {

	slice := make([]int, size, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Intn(999) - rand.Intn(999)
	}
	return slice
}

func main() {
	version := "v1"
	if len(os.Args) == 2 {
		version = os.Args[1]
	}

	f, err := os.OpenFile(version+".trace", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	trace.Start(f)
	defer trace.Stop()

	for i := 0; i < 10000; i++ {
		s := generateSlice(10)

		switch version {
		case "v1":
			mergesortv1(s)
		case "v2":
			mergesortv2(s)
		case "v3":
			mergesortv3(s)
		}
	}

}
