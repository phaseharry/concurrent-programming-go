package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
using atomic operations instead of manually using a mutex to protect shared resource variables (only available for certain operations)
only use atomic operations for preventing race conditions of multiple goroutines accessing the same shared resources as there is a performance
penalty compared to regular operations. this happens because when atomic operations are used, a lot of compiler and system optimizations have to be forfieted.
ex, when a variable is accessed repeatedly the system will keep that variable in the processor cache, making access to that variable faster.
it might periodically flush the variable back to main memory if it is running out of cache space in the processor.
when using atomics, the system needs to ensure other executions running in parallel get the latest version of the variable. It does this by
ensuring the processor cached values are always in sync by flushing the processor cached value to main memory and invalidating other execution's processor
cached value. This requirement to keep caches consistent with each other reduces performance so atomic operations should only be used when actually required.

also using .LoadInt32 is not strictly required at the end of the main function as we wait until our 2 goroutines are done running before we
actually read and print it so we know for a fact there is no other thread of execution making changes to it, but using it is best practice for any
resource that is shared between threads of execution. This insures we get the latest value from main memory and not an outdated cache value.
*/
func main() {
	money := int32(100)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		stingy(&money)
		wg.Done()
	}()
	go func() {
		spendy(&money)
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("Money in account : ", atomic.LoadInt32(&money))
}

func stingy(money *int32) {
	for range 1_000_000 {
		atomic.AddInt32(money, 10)
	}
	fmt.Println("Stingy done")
}

func spendy(money *int32) {
	for range 1_000_000 {
		atomic.AddInt32(money, -10)
	}
	fmt.Println("Spendy done")
}
