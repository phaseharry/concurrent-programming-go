package main

import "fmt"

/*
demo of go channels being first-class objects (can be passed into functions / created programmatically during runtime)
iterate through numbers in ascending order (2 -> 99,999).
*/
func main() {
	/*
		initial prime channel (2) since 2 will be the first number getting passed into this channel.
		then it consumes all the numbers after 2 and determine if those numbers are prime or not. If not
		prime and divisible by the current multipleFilter then it will be tossed out. If it's not divisible
		by current multipleFilter, send the number to the next multipleFilter for checking. If the next multipleFilter
		does not exist, then create a new multipleFilter goroutine and it will divide by that number and be that multipleFilterer.
	*/
	numbers := make(chan int)
	quit := make(chan int)
	go primeMultipleFilter(numbers, quit)

	// sending 2 -> 100,000 into the primeNumberFilter
	for i := 2; i < 100_000; i++ {
		numbers <- i
	}
	close(numbers)

	// blocking main goroutine until a quit signal/message has been received
	<-quit
}

/*
2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16

2 is the first number that's passed into the first primeMultipleFilter goroutine so that first goroutine will be the 2 prime number filterer.
This new goroutine will check all numbers after 2 and filter out numbers that are divisible by 2. For numbers not divisible by 2, it will pass it
to the next primeMultipleFilter goroutine in the chain (3) and have that check if it's divisible by 3. If so, it will be tossed out otherwise, it will
do the same and pass it to the next goroutine in the chain.
*/
func primeMultipleFilter(numbers <-chan int, quit chan<- int) {
	var right chan int
	p := <-numbers
	fmt.Println(p)

	for n := range numbers {
		/*
			discarding numbers that a multiples of p as p is the current number.
			if the prime number directly after p has not be encountered yet, create
			a new channel and spawn a new goroutine that will filter out numbers that
			are multiples of the next prime number.
		*/
		if n%p != 0 {
			if right == nil {
				right = make(chan int)
				go primeMultipleFilter(right, quit)
			}
			right <- n
		}
	}
	if right == nil {
		close(quit)
	} else {
		close(right)
	}
}
