package main

import "fmt"

/*
Quitting Channels

pattern of using a channel to signal the termination of a
goroutine once we're done processing something
*/
func main() {
	numbers := make(chan int)
	quit := make(chan int)
	printNumbers(numbers, quit)
	next := 0

	// iterating infinite times with no number here being a signal to stop iterating
	for i := 1; ; i++ {
		next += i
		select {
		/*
		   sending numbers through the numbers channel to be printed.
		   once the goroutine in printNumbers finish consuming 10 numbers,
		   it will close the quit channel. At that point, the quit channel
		   will start having the default values of type int being consumed
		   and that will be our signal that we've consumed the right amount of
		   numbers and can terminate the main goroutine
		*/
		case numbers <- next:
		case <-quit:
			fmt.Println("Quitting number generation")
			return
		}
	}
}

func printNumbers(numbers <-chan int, quit chan int) {
	go func() {
		/*
		   consumes 10 messages from numbers channel, then
		   close the quit channel to signal that the main goroutine
		   can be unblocked and stop sending messages to the numbers channel
		*/
		for range 10 {
			fmt.Println(<-numbers)
		}
		close(quit)
	}()
}
