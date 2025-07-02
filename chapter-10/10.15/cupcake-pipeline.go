package main

import (
	"fmt"

	listing10_13 "github.com/phaseharry/concurrent-programming-go/chapter-10/10.13"
)

func main() {
	input := make(chan int)
	quit := make(chan int)

	/*
	   Demo of pipelining in which there's multiple threads of execution (goroutines / person working in cupcake making process)
	   work together to process a job. Each goroutine works on a different task within the process and only on that. Once
	   that goroutine finishes their work, it pipes it into the next goroutine in the pipe and it does the same.
	   This increases throughput over sequential processing as the 10 batches of cupcake will be done processing faster.
	   Each person will only do one tasks and pass their results onto the next person in the process.

	   run: time go run cupcake-pipeline.go

	   Notes regarding speeding up a pipeline process:
	   - to increase throughput (number of jobs completed within X unit of time), you have to focus on the bottleneck
	   and decrease the time the slowest process takes
	   - to increase latency (the time it takes for a single job to run), you just need to decrease the time it takes for
	   any / all tasks of the job
	*/
	finishedPreparingTray := AddOnPipe(quit, listing10_13.PrepareTray, input)
	finishedMixture := AddOnPipe(quit, listing10_13.Mixture, finishedPreparingTray)
	finishedBaking := AddOnPipe(quit, listing10_13.Bake, finishedMixture)
	finishedAddingToppings := AddOnPipe(quit, listing10_13.AddToppings, finishedBaking)
	finishedBoxing := AddOnPipe(quit, listing10_13.Box, finishedAddingToppings)

	/*
	   create a new goroutine that sends 10 integers into the initial input queue
	   indicating the batch number.
	*/
	requestedBatches := 10

	go func() {
		for i := range requestedBatches {
			input <- i
		}
	}()

	// blocking until we've finished all 10 batches of cupcake
	for range requestedBatches {
		fmt.Println(<-finishedBoxing, "received")
	}
}

/*
common util function that takes in a channel and consumes messages from it
and pipe the message in a newly created channel.
uses the quit channel pattern to close the created output channel to stop
further processing.
takes in a function that processes the messages from the incoming channel
and return its output as input to the next channel in the pipeline
*/
func AddOnPipe[X, Y any](q <-chan int, f func(X) Y, in <-chan X) chan Y {
	output := make(chan Y)

	go func() {
		defer close(output)
		for {
			select {
			case <-q:
				return
			case input := <-in:
				output <- f(input)
			}
		}
	}()

	return output
}
