package wsemaphore

import "sync"

type WeightedSemaphore struct {
	cond    *sync.Cond
	permits int
}

func NewWeightedSemaphore(permits int) *WeightedSemaphore {
	return &WeightedSemaphore{
		cond:    sync.NewCond(&sync.Mutex{}),
		permits: permits,
	}
}

func (ws *WeightedSemaphore) Acquire(permitsRequired int) {
	ws.cond.L.Lock()

	/*
	   Since we're doing a weighted semaphore that gets passed in
	   the value of the permit a goroutine wants to use, we have to
	   check that the semaphore has the amount of permits to give.
	   If not, the goroutine will Wait() until the amount is available.
	*/
	for permitsRequired > ws.permits {
		ws.cond.Wait()
	}

	ws.permits -= permitsRequired

	ws.cond.L.Unlock()
}

func (ws *WeightedSemaphore) Release(permitsReleased int) {
	ws.cond.L.Lock()

	ws.permits += permitsReleased

	/*
	   using Broadcast() so all blocked goroutines can attempt to acquire permits.
	   since this is weighted, the current goroutine may have 5 permits and releasing
	   them can allow all 5 goroutines that all have 1 permit required to process.
	   for a normal semaphore where each goroutine only requires 1 permit, we used .Signal()
	   to only wake one goroutine that's blocked since each goroutine only needed one to run.
	   since each goroutine in a weighted semaphore has a variable permit requirement,
	   we have to wake all existing blocked goroutines because there's a change all can process
	   due to the high weighted goroutine releasing its permits.
	*/
	ws.cond.Broadcast()

	ws.cond.L.Unlock()
}
