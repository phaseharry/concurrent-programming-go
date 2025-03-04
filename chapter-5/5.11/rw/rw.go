package rw

import "sync"

/*
The ReadWriteMutex we implemented in 4.14 is a read preferred
readers-writer mutex. There are scenarios in which the writer
can be starved and never be able to be allowed to acquire the lock
to make a write request if there a constant read requests and the
read counter will never hit zero. This example will demonstrate a
write-preferring readers-writer mutex in which many goroutines can
acquire a read lock, but once a write lock gets called, additional
reads will be blocked until all writes resolve and release the lock
back.
*/

/*
readersCounter - used to track current goroutines that have the read lock
writersWaiting - used to track current goroutines that wants to acquire the writer lock
writerActive - bool to indicate whether or not a goroutine has a write lock or not
cond - the condition we used to synchronize everything.

in this implementation,
1. if there is a goroutine attempting to acquire a write lock and there are outstanding read locks
  - any outstanding read goroutines will be allowed to finish (readersCount -> 0)
  - any new read locks will be blocked and must wait until this goroutine acquire the lock and release it when it's done before a read lock can be acquired (

2. if there is a goroutine attempting to acquire a write lock and there is an outstanding write lock (writerActive == true)
  - this goroutine must wait and we track that there's a write goroutine waiting (increment writersWaiting)
  - any goroutines that tries to acquire the read lock must wait for both the existing write goroutine to finish & the writers that are waiting to process before they can acquire the lock

3. if there are existing read goroutines that have read locks and a writer goroutine tries to acquire a writer lock
  - that writer goroutine will be blocked until all existing read goroutines release the lock (readersCount -> 0)
*/
type ReadWriteMutex struct {
	readersCounter int
	writersWaiting int
	writerActive   bool
	cond           *sync.Cond
}

func NewReadWriteMutex() *ReadWriteMutex {
	return &ReadWriteMutex{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (rw *ReadWriteMutex) ReadLock() {
	// acquire lock to make changes within the ReadWriteMutex
	rw.cond.L.Lock()

	/*
	   if there's existing goroutines that are waiting for the write lock
	   or there's a goroutine that's active writing, then block the goroutine
	   that attempted the read lock and have it wait until all writes are finished.
	*/
	for rw.writersWaiting > 0 || rw.writerActive {
		rw.cond.Wait()
	}

	/*
	   once we're able to acquire the read lock once all writes are done or not writes at all, then
	   keep track of all existing goroutines that have an outstanding read lock.
	*/
	rw.readersCounter++

	// unlock so another goroutine can attempt to get a ReadLock
	rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) WriteLock() {
	rw.cond.L.Lock()

	/*
		Increment the writersWaiting, assuming that the current goroutine might have to wait.
		If the `writerActive` flag is false then it doesn't have to wait and the writersWaiting
		value will be decremented to zero and it has the writeLock.
		Otherwise, if there's an existing goroutine with the writeLock then block the current goroutin
		trying to acquire the writeLock.
	*/
	rw.writersWaiting++

	for rw.readersCounter > 0 || rw.writerActive {
		rw.cond.Wait()
	}

	rw.writersWaiting--
	rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) ReadUnlock() {
	rw.cond.L.Lock()

	/*
	   unlocking will decrement the `readersCounter` value.
	   if it reaches zero then that means the current goroutine is the last goroutine that had a read lock,
	   so we need to broadcast so other goroutines that are attempting to get the write lock
	   can be unblocked, acquire the condition's lock and attempt to get the write lock.
	   (saying attempt because there could be multiple goroutines trying to acquire the write lock.)
	   if there are goroutines attempting to get read locks then those will be blocked until all
	   writer locks have been released.
	*/
	rw.readersCounter--

	if rw.readersCounter == 0 {
		rw.cond.Broadcast()
	}

	rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) WriteUnlock() {
	rw.cond.L.Lock()

	// setting writer active to false & broadcasting so other write & read goroutines can get lock and process
	rw.writerActive = false

	rw.cond.Broadcast()

	rw.cond.L.Unlock()
}
