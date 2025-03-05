package semaphore

import "sync"

/*
implementation of a semaphore using conditions and mutuxes.
go doesn't ship an implementation as part of standard library, but they do have an implementation
that can be used.

semaphores allow N executions to operation concurrently as oppose to just a single one execution
getting the mutual exclusive access and other executions are blocked until that one execution is done.

- mutex allows a single goroutine access at a time
- semaphores allow N goroutines access at a time

semaphores do it by having a permit field that holds the max number of goroutines that can be running
at any instance of time.
note: a semphore with 1 permit is effectively a mutex since it allows only one goroutine to run at a time.
*/

/*
- takes in the an int value to determine the number of concurrent gorountines
allowed to execute.
- a condition to get exclusive access to updating the permits and blocking other
processes when we're out of permits to give
*/
type Semaphore struct {
	permits int
	cond    *sync.Cond
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		permits: n,
		cond:    sync.NewCond(&sync.Mutex{}),
	}
}

func (rw *Semaphore) Acquire() {
	rw.cond.L.Lock() // the condition will acquire the lock

	/*
	   if a goroutine is attempting to acquire access and we're out of permits
	   then the goroutine that's attemtping to acquire access will have to wait.
	*/
	for rw.permits <= 0 {
		/*
		   the goroutine that's waiting will release the condition lock so another
		   goroutine can acquire access. once that current goroutine get a signal,
		   it will acquire the condition lock again and run through the loop again.
		   if there's permits > 0, then we will escape the for loop and not .Wait().
		   we decrement permits since the current goroutine will have acquired access
		   to the semaphore and release the condition lock so another goroutine can
		   attempt to Acquire or Release their access to semaphore.
		*/
		rw.cond.Wait()
	}
	rw.permits--
	rw.cond.L.Unlock()
}

func (rw *Semaphore) Release() {
	rw.cond.L.Lock()
	/*
	   when a goroutine releases its access to the semaphore, it just increment
	   the permits and call .Signal(). Since all of our .Waits() within the
	   semaphore are goroutines attempting to acquire access, we call .Signal()
	   and let a random goroutine have access to it since we've only incremented
	   the permit by one.
	*/
	rw.permits++

	rw.cond.Signal()
	rw.cond.L.Unlock()
}
