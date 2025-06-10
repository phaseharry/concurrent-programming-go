package exercise_9_3_4

import "fmt"

func Drain[T any](quit <-chan int, input <-chan T) {

	go func() {
		moreData := true
		for moreData {
			select {
			case message, moreData := <-input:
				if moreData {
					fmt.Println("throwing away message:", message)
				}
			case <-quit:
				return
			}
		}
	}()

}
