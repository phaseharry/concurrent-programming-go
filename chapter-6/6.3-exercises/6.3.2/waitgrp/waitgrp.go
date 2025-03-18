package waitgrp

import "sync"

/*
new implementation of WaitGrp that allows for adding to
the wait counter after it has been initialized.
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

func (wg *WaitGrp) Add(delta int) {
	wg.cond.L.Lock()
	wg.groupSize += delta
	wg.cond.L.Unlock()
}

func (wg *WaitGrp) Wait() {
	wg.cond.L.Lock()

	for wg.groupSize > 0 {
		wg.cond.Wait()
	}

	wg.cond.L.Unlock()
}

func (wg *WaitGrp) Done() {
	wg.cond.L.Lock()
	wg.groupSize -= 1
	if wg.groupSize == 0 {
		wg.cond.Broadcast()
	}
	wg.cond.L.Unlock()
}

/*
Non-blocking TryWait() similiar to mutex's TryLock().
Returns true if Wait() is ready otherwise false.
This will not block the goroutine that calls it
*/
func (wg *WaitGrp) TryWait() bool {
	wg.cond.L.Lock()
	waitGroupEmpty := wg.groupSize == 0
	wg.cond.L.Unlock()
	return waitGroupEmpty
}
