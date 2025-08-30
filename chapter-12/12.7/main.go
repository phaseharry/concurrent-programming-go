package main

import (
	"fmt"
	"sync/atomic"
)

/*
demo of atomic.CompareAndSwap
*/

func main() {
	/*
	   CompareAndSwap is an atomic operations so if it's called it will not be interrupted by
	   other executions that's calling CompareAndSwap on the same pointer variable / memory.
	   It takes in 3 parameters. The first parameter is the pointer to variable that will be modified
	   by multiple executions, the second parameter is the value that is checked against and the third
	   parameter is the value that gets assigned to that pointer location if the initial value at that
	   location is equal to the second parameter.
	*/
	number := int32(17)
	result := atomic.CompareAndSwapInt32(&number, 17, 19)
	/*
		since the value in the &number location is 17 and it matches the 2nd paramter, 19
		will be atomically swapped into the &number location
	*/
	fmt.Printf("17 <- swap(17, 19): result %t, value: %d\n", result, number)

	number = int32(23)
	result = atomic.CompareAndSwapInt32(&number, 17, 19)
	/*
		because number was reassigned to 23, it will not match 17 and be swapped with 19 again. Instead
		the &number address will still hold the value of 23.
	*/
	fmt.Printf("23 <- swap(17, 19): result %t, value: %d\n", result, number)
}
