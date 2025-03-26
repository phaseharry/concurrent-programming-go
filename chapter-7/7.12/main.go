package main

import "fmt"

func main() {
	resultCh := make(chan []int)
	go func() {
		/*
		   create a goroutine to call findFactors of 3419110721
		   and send that result to the resultCh.
		   while we do this, we will have the main goroutine
		   call the findFactors of 4033836233.
		   if we sequentially called for both, then we wouldn't be able
		   to take advantage of multiple cores if they were available.
		   this way, there's 2 goroutines processing concurrently.
		*/
		resultCh <- findFactors(3419110721)
	}()
	// main goroutine will block until factors of 4033836233 are calculated
	fmt.Println(findFactors(4033836233))
	/*
	   will block until the child goroutine finishes
	   calculating factors of 3419110721 since we're
	   waiting to consume a message from an unbuffered
	   channel
	*/
	fmt.Println(<-resultCh)
}

// find all factors for a given number
func findFactors(number int) []int {
	result := make([]int, 0)
	for i := 1; i <= number; i++ {
		if number%i == 0 {
			result = append(result, i)
		}
	}
	return result
}
