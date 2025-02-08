package main

import (
	"fmt"
	"sync"
)

func main() {
	rwMutex := ReadWriteMutex{}
	fmt.Println("Acquiring Readlock")
	rwMutex.ReadLock()
	fmt.Println("Acquiring Readlock again")
	rwMutex.ReadLock()
	fmt.Println("Trying Readlock", rwMutex.TryReadLock())
	fmt.Println("Trying Writelock", rwMutex.TryWriteLock())
	rwMutex.ReadUnlock()
	rwMutex.ReadUnlock()
	rwMutex.ReadUnlock()
	fmt.Println("Trying Writelock", rwMutex.TryWriteLock())
	fmt.Println("Trying Readlock", rwMutex.TryReadLock())
}

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

// add try write lock for exercise 2
func (rw *ReadWriteMutex) TryWriteLock() bool {
	return rw.globalLock.TryLock()
}

// add try read lock for exercise 3
func (rw *ReadWriteMutex) TryReadLock() bool {
	acquiredReadLock := rw.readersLock.TryLock()

	/*
	 wasn't able to acquire read lock, meaning another gorountine
	 has acquired the lock to increment the readerCounter
	*/
	if !acquiredReadLock {
		return false
	}

	acquiredGlobalLock := true
	/*
		if this is the first readersLock we're trying to acquire, we need to be able
		to get the global lock as well so a goroutine cannot write to the shared memory space
		the same time as we read it
	*/
	if rw.readersCounter == 0 {
		acquiredGlobalLock = rw.TryWriteLock()
	}

	if acquiredGlobalLock {
		rw.readersCounter += 1
	}

	rw.readersLock.Unlock()

	/*
	 if we were able to acquire the globalLock as well as the initial readLock (to increment the readers counter)
	 then we know that we were able to get the readLock
	*/
	return acquiredGlobalLock
}

func (rw *ReadWriteMutex) WriteUnlock() {
	rw.globalLock.Unlock()
}
