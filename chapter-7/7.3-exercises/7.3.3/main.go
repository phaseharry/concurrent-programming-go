package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(10)

	for range 10 {
		go findFactors(randRange(100000, 500000), &wg)
	}

	wg.Wait()
}

func findFactors(number int, wg *sync.WaitGroup) {
	result := make([]int, 0)
	for i := 1; i <= number; i++ {
		if number%i == 0 {
			result = append(result, i)
		}
	}
	fmt.Printf("%d's factors: %v\n", number, result)
	wg.Done()
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
