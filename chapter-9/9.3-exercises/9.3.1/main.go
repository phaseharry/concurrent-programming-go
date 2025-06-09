package main

import (
	"fmt"
)

func main() {
	quit := make(chan int)
	defer close(quit)

	squares := GenerateSquares(quit)
	for val := range squares {
		fmt.Println(val)
	}
}

func GenerateSquares(quit <-chan int) <-chan int {
	squaresChannel := make(chan int)

	go func() {
		defer close(squaresChannel)
		currentNumber := 1
		for {
			select {
			case squaresChannel <- (currentNumber * currentNumber):
				currentNumber += 1
			case <-quit:
				return
			}
		}
	}()

	return squaresChannel
}
