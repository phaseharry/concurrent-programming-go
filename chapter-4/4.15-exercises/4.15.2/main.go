package main

import "sync"

type ReadWriteMutex struct {
	readersCounter int // keeps track of goroutines that have currently acquired a Read lock

	/*
		mutex for synchronizing readers access.
		to used call Lock when we increment & decrement the readersCounter to prevent
		a race condition when that's getting updated
	*/
	readersLock sync.Mutex

	globalLock sync.Mutex // mutex for blocking any writers access
}

func (rw *ReadWriteMutex) ReadLock() {
	rw.readersLock.Lock()
	rw.readersCounter++

	/*
		if this current ReadLock call is the first one, we will only have a readersCounter of 1
		so we also need to lock the globalLock to prevent any Writes from updating critical space
	*/
	if rw.readersCounter == 1 {
		rw.globalLock.Lock()
	}

	// call unlock once the readersCounter has been incrementedx
	rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) WriteLock() {
	/*
		if the globalLock has already been acquired, then we need to wait for it to be
		unlocked first. The goroutine that calls this will be blocked at the call
	*/
	rw.globalLock.Lock()
}

func (rw *ReadWriteMutex) ReadUnlock() {
	rw.readersLock.Lock()
	rw.readersCounter--

	/*
		If there are no more goroutines with readers lock then we can
		unlock the global lock so any blocked write call can be unblocked and gain
		access to the shared memory space
	*/
	if rw.readersCounter == 0 {
		rw.globalLock.Unlock()
	}

	rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) TryWriteLock() bool {
	return rw.globalLock.TryLock()
}

func (rw *ReadWriteMutex) WriteUnlock() {
	rw.globalLock.Unlock()
}
