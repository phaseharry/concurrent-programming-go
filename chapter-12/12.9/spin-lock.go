package main

import (
	"runtime"
	"sync"
	"sync/atomic"
)

/*
Demo of a SpinLock (mutex)

SpinLock is a type of lock in which an execution will go into a loop to try to get a hold
of a lock repeatedly until the lock becomes available.

a value of 0 indicates that the lock is free to be acquired by a goroutine
while a value of 1 indicates that another goroutine already acquired exclusive access

this implements Go's Locker interface so it can be used in place of the built-in sync.Mutex

issues with current implementation:
- occurs when there is high resource contention (ex. when a goroutine is hogging a lock for a long time),
other goroutines/execution will bew wasting valuable CPU cycles while spinning(looping repeated) and waiting for the lock to be released.
- in our case, goroutines will be stuck in a loop, calling CompareAndSwap repeatedly until a goroutine calls Unlock to free
the resource. This looping and waiting wastes valuable CPU time that could be used to execute other tasks
*/
type SpinLock int32

func (s *SpinLock) Lock() {
	/*
	   atomically compare the current value of the SpinLock. If the attempted CompareAndSwap
	   returns true when compared with 0 then that means the lock was free to lock it with 1.
	   so if it returned false, meaning that compared with 0, it did not match indicating that the
	   lock was already taken so keep looping until lock can be acquired.
	   During each loop iteration, call the Go scheduler to give execution time to other goroutines.
	*/
	for !atomic.CompareAndSwapInt32((*int32)(s), 0, 1) {
		runtime.Gosched()
		/*
			here, we are yielding the execution to allow another execution to run. The issue is that the runtime or operating system
			does not that the current execution/goroutine is waitng for a lock to become available so it is highly likely that the current
			execution will run more times before another execution can use CPU time. to help with this, OS provides a concept known as a futex.

		*/
	}
}

func (s *SpinLock) Unlock() {
	/*
	   atomically store 0 to the SpinLock value to allow another goroutine attempt
	   to acquire exclusive access to the spin lock
	*/
	atomic.StoreInt32((*int32)(s), 0)
}

func NewSpinLock() sync.Locker {
	var lock SpinLock
	return &lock
}
