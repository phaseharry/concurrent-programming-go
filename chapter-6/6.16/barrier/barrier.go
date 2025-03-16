package barrier

import "sync"

type Barrier struct {
	size      int
	waitCount int
	cond      *sync.Cond
}

func NewBarrier(size int) *Barrier {
	condVar := sync.NewCond(&sync.Mutex{})
	/*
	   initializing the Barrier to have a current waitCount of 0
	   and a size of size. Once waitCount hits the size value,
	   all goroutines that have called .Wait() will be unblocked
	   as all those goroutines have processed up to the point we
	   allowed them to before we block them and synchronize their
	   continued processing.
	*/
	return &Barrier{size, 0, condVar}
}

func (b *Barrier) Wait() {
	b.cond.L.Lock()
	/*
	   incrementing the waitCount by 1 for the goroutine
	   that calls .Wait(). That goroutine will be blocked
	   until the Barrier's waitSize has reached size through
	   other goroutines calling .Wait() and incrementing waitCount.
	   Once the waitCount has reached the Barrier size value, we
	   reset the Barrier count so it can be reused and call .Broadcast()
	   so all goroutines that were blocked by the Barrier can continue processing.
	*/
	b.waitCount += 1

	if b.waitCount == b.size {
		b.waitCount = 0
		b.cond.Broadcast()
	} else {
		b.cond.Wait()
	}
	b.cond.L.Unlock()
}
