package main

import (
	"fmt"
	"math"
	"time"
)

const (
	PASSWORD_TO_GUESS = "go far"
	ALPHABET          = " abcdefghijklmnopqrstuvwxyz"
)

func main() {
	finished := make(chan int)
	passwordFound := make(chan string)

	/*
	   going through all the possible password combinations (1 -> 387,420,488)
	   and split guesses to multiple child goroutines. Each goroutine will make
	   guesses for 10,000,000 possible passwords.
	   ex.
	   - goroutine 1 -> will make a guess for passwords 1 -> 10,000,000
	   - goroutine 2 -> 10,000,000 -> 20,000,000
	   - goroutine 3 -> 20,000,000 -> 30,000,000
	   - etc.
	*/
	for i := 1; i <= 387_420_488; i += 10_000_000 {
		go guessPassword(
			i,
			int(math.Min(float64(i+10_000_000), float64(387_420_488))),
			finished,
			passwordFound,
		)
	}
	/*
	   blocking the main goroutine until one of the separate child goroutines is
	   able to guess the correct password.
	*/
	fmt.Println("password found:", <-passwordFound)
	close(passwordFound)
	time.Sleep(5 * time.Second)
}

func guessPassword(start int, upto int, stop chan int, result chan string) {
	for currentNum := start; currentNum < upto; currentNum += 1 {
		select {
		/*
		   will continue guessing the password until we either get the password
		   or we run out of numbers to guess within the current range of [start, upto)

		   normally there will never be a message in the stop channel since we're not
		   sending any messages to that channel so the default case will always run.
		   if we're able to guess the correct password then we will send the correct
		   number to the result channel to unblock the main goroutine.
		   we will also close the stop channel.
		   this will cause any goroutines that consumes messages from the stop channel as part of its select case to start
		   consuming the default int value. when that happens we know we've found the password in
		   another goroutine so we will return and terminate the other goroutines.
		*/
		case <-stop:
			fmt.Printf("Stopped at %d [%d, %d)\n", currentNum, start, upto)
			return
		default:
			if toBase27(currentNum) == PASSWORD_TO_GUESS {
				fmt.Printf("Found at %d [%d, %d)\n", currentNum, start, upto)
				result <- toBase27(currentNum)
				close(stop)
				return
			}
		}
	}
}

/*
not important, but just converts an int to a string assuming it's a base27 string (see ALPHABET, 27 possible characters).
the possible strings for password is between the lengths of 1 and 6 with those possible characters.
possbile strings goes from "a" -> "zzzzzz" (1 -> 387,420,488)
1 -> 27^6 - 1
*/
func toBase27(n int) string {
	result := ""
	for n > 0 {
		result = string(ALPHABET[n%27]) + result
		n /= 27
	}
	return result
}
