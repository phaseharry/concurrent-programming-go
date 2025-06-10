package main

import (
	exercise_9_3_1 "github.com/phaseharry/concurrent-programming-go/chapter-9/9.3-exercises/9.3.1"
	exercise_9_3_2 "github.com/phaseharry/concurrent-programming-go/chapter-9/9.3-exercises/9.3.2"
	exercise_9_3_3 "github.com/phaseharry/concurrent-programming-go/chapter-9/9.3-exercises/9.3.3"
	exercise_9_3_4 "github.com/phaseharry/concurrent-programming-go/chapter-9/9.3-exercises/9.3.4"
)

func main() {
	quitChannel := make(chan int)

	exercise_9_3_4.Drain(
		quitChannel,
		exercise_9_3_3.Print(
			quitChannel,
			exercise_9_3_2.TakeUntil(
				func(s int) bool { return s <= 10000 },
				quitChannel,
				exercise_9_3_1.GenerateSquares(quitChannel),
			),
		),
	)

	<-quitChannel
}
