package main

import (
	"fmt"
	"math"
	"math/rand"
)

/*
example showing the select's statement having both cases
from sending messages to a channel and a case for consuming
messages from a channel as well.
*/
func main() {
	numbersChannel := make(chan int)
	primes := primesOnly(numbersChannel)
	for i := 0; i < 100; {
		select {
		/*
		   in this case, we win send random numbers through the numbersChannel
		   and the primesOnly function has a goroutine running that consumes messages
		   from that channel.
		   the primesOnly goroutine will determine if that number is a primeNumber and if it
		   is, send that number through the primes channel for the main goroutine to process
		   by consuming it in the other select case.
		   we will do this until we have found 100 random numbers that are prime
		*/
		case numbersChannel <- rand.Intn(1000000000) + 1:
		case p := <-primes:
			fmt.Println("Found prime:", p)
			i++
		}
	}
}

func primesOnly(inputs <-chan int) <-chan int {
	results := make(chan int)
	/*
	   create a chan and return that channel so another goroutine
	   can consume its messages.
	   Also start a goroutine in which it consumes random numbers that
	   gets send through the inputs channel
	*/
	go func() {
		for c := range inputs {
			isPrime := c != 1
			for i := 2; i <= int(math.Sqrt(float64(c))); i++ {
				if c%2 == 0 {
					isPrime = false
					break
				}
			}
			if isPrime {
				results <- c
			}
		}
	}()
	return results
}
