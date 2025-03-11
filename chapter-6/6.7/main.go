package main

import "sync"

/*
new implementation of WaitGrp that allows for Adding to
the wait counter after it has been initialized
*/
type WaitGrp struct {
	groupSize int
	cond      *sync.Cond
}

func NewWaitGrp() *WaitGrp {
	return &WaitGrp{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

/*
acquire lock on the condition's mutuex to ensure
we are safely updating the groupSize. unlock once
we're done incrementing the groupSize by delta
*/
func (wg *WaitGrp) Add(delta int) {
	wg.cond.L.Lock()
	wg.groupSize += delta
	wg.cond.L.Unlock()
}

/*
the goroutine that calls .Wait() will be
blocked while the groupSize is greater than 0.
only when it gets a signal and it checks that
groupSize == 0 will it be able to be unblocked
and process.
*/
func (wg *WaitGrp) Wait() {
	wg.cond.L.Lock()

	for wg.groupSize > 0 {
		wg.cond.Wait()
	}

	wg.cond.L.Unlock()
}

/*
when a gorountine calls .Done(), we decrement the groupSize
and if the groupSize == 0 then we Broadcast() to let all
blocked goroutines that called .Wait() that they can become
unblocked and continue processing
*/
func (wg *WaitGrp) Done() {
	wg.cond.L.Lock()
	wg.groupSize -= 1
	if wg.groupSize == 0 {
		wg.cond.Broadcast()
	}
	wg.cond.L.Unlock()
}
