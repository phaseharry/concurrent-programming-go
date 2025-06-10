package exercise_9_3_2

import (
	"fmt"
)

func TakeUntil[K any](f func(K) bool, quit chan int, input <-chan K) <-chan K {
	output := make(chan K)

	go func() {
		defer close(output)

		moreData := true
		continueProcessing := true

		for moreData && continueProcessing {
			select {
			case message, moreData := <-input:
				if moreData {
					continueProcessing = f(message)
					if continueProcessing {
						output <- message
					}
				}
			case <-quit:
				return
			}
		}

		// close quit channel as well if we're not processing anymore
		fmt.Println("continueProccessing:", continueProcessing)
		if !continueProcessing {
			close(quit)
		}
	}()

	return output
}
