package exercise_9_3_3

import "fmt"

func Print[T any](quit <-chan int, input <-chan T) <-chan T {
	output := make(chan T)

	go func() {
		defer close(output)

		moreData := true

		for moreData {
			select {
			case message, moreData := <-input:
				if moreData {
					fmt.Println(message)
					output <- message
				}
			case <-quit:
				return
			}
		}
	}()

	return output
}
