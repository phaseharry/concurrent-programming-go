package waitgrp

import (
	"6.3/semaphore"
)

/*
implementing a WaitGroup using the Semaphore we created in chapter 5
using WaitGrp to not conflict with go's WaitGroup implementation
note: this WaitGroup implementation does allow Adding new values to increase
the wait parameter. that is supported in go's implementation. will be created in 6.5
*/
type WaitGrp struct {
	sema *semaphore.Semaphore
}

func NewWaitGrp(size int) *WaitGrp {
	/*
	   initializing the semaphore with 1 - size so that it is immediately blocked
	   and calling Done() will call the semaphore's Release() to increment the
	   permit count and once that permit count is eventually 1, the WaitGrp's .Wait()
	   call can finally acquire the permit on the semaphore and process the goroutine
	*/
	return &WaitGrp{
		sema: semaphore.NewSemaphore(1 - size),
	}
}

func (wg *WaitGrp) Wait() {
	wg.sema.Acquire()
}

func (wg *WaitGrp) Done() {
	wg.sema.Release()
}
