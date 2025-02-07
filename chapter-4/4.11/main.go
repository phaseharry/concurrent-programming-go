package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {
	/*
		Readers-Writer mutex
	*/
	mutex := sync.RWMutex{}

	// set first 10000 events
	var matchEvents = make([]string, 0, 10000)
	for i := 0; i < 10000; i++ {
		matchEvents = append(matchEvents, "match events")
	}

	go matchRecorder(&matchEvents, &mutex)

	start := time.Now()

	// handle 5000 clients pulling events
	for j := 0; j < 5000; j++ {
		go clientHandler(&matchEvents, &mutex, start)
	}

	time.Sleep(100 * time.Second)
}

func matchRecorder(matchEvents *[]string, mutex *sync.RWMutex) {
	for i := 0; ; i++ { // infinite loop to append events to the match events list even 200 milliseconds
		/*
			Using the Write lock of the ReadWrite lock so if it's called, it will be blocked if there's
			any existing write or read lock outstanding. The goroutine will only be able to acquire the Write lock
			if there's no Read or Write lock outstanding. This will prevent race conditions when we make update to
			the shared memory space of "matchEvents".

			writes an event to the match event every 200 milliseconds
		*/
		mutex.Lock()
		*matchEvents = append(*matchEvents, "Match event "+strconv.Itoa(i))
		mutex.Unlock()
		time.Sleep(200 * time.Millisecond)
		fmt.Println("Appended match event")
	}
}

func clientHandler(matchEvents *[]string, mutex *sync.RWMutex, startTime time.Time) {
	for i := 0; i < 100; i++ {
		/*
			when a goroutine tries to acquire a Read lock and there's only other Read locks acquired, then
			the current goroutine will not be blocked and can acquire the Read lock and gain access to the
			shared memory space. We still need to call Read Unlock to track that the memory space isn't being used anymore,
			so Write locks can be acquired and make changes to that memory space when it's needed and there's no outstanding
			goroutine with a Read or Write lock. However if there's an outstanding Write lock, then the goroutine will be blocked
			and will have to wait until the goroutine finishes its write.
		*/
		mutex.RLock()
		allEvents := copyAllEvents(matchEvents)
		mutex.RUnlock()
		timeTaken := time.Since(startTime)
		fmt.Println(len(allEvents), "events copied in", timeTaken)
	}
}

func copyAllEvents(matchEvents *[]string) []string {
	allEvents := make([]string, 0, len(*matchEvents))
	for _, e := range *matchEvents {
		allEvents = append(allEvents, e)
	}
	return allEvents
}
