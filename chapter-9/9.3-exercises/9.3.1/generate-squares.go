package exercise_9_3_1

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
