package main

func TakeUntil[K any](f func(K) bool, quit chan int, input <-chan K) <-chan K {
	output := make(chan K)

	go func() {
		defer close(output)

		moreData := true
		continueProcessing := true

		for moreData {
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
		if !continueProcessing {
			close(quit)
		}
	}()

	return output
}
