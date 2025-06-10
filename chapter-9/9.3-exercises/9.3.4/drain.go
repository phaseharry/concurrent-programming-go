package exercise_9_3_4

func Drain[T any](quit <-chan int, input <-chan T) {

	go func() {
		moreData := true
		for moreData {
			select {
			case <-input:
			case <-quit:
				return
			}
		}
	}()

}
