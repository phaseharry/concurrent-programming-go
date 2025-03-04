package main

import (
	"5.11/rw"
	"fmt"
	"time"
)

func main() {
	rwMutex := rw.NewReadWriteMutex()

	for range 2 {
		go func() {
			for {
				rwMutex.ReadLock()
				time.Sleep(1 * time.Second)
				fmt.Println("Read done")
				rwMutex.ReadUnlock()
			}
		}()
	}
	time.Sleep(1 * time.Second)
	rwMutex.WriteLock()
	fmt.Println("Write finished")
}
