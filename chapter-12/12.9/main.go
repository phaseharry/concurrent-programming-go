package main

import (
	"sort"
	"sync"
)

// using a built in sync.Locker
type Flight struct {
	Origin, Dest string
	SeatsLeft    int
	Locker       sync.Locker
}

func Book(flights []*Flight, seatsTobook int) bool {
	bookable := true

	/*
		sorts flights in alphabetical order based on their origin and destination names.
		this will ensure no goroutines can only acquire partial locks but be blocked by other
		another goroutine who has the exclusive access to the remaining and be blocked as well.
		prevent deadlocks from occurring
	*/
	sort.Slice(flights, func(a, b int) bool {
		flightA := flights[a].Origin + flights[a].Dest
		flightB := flights[b].Origin + flights[b].Dest
		return flightA < flightB
	})

	for _, f := range flights {
		f.Locker.Lock()
	}

	/*
		once all locks were acquired by the current goroutine, check if all requested
		flights a bookable. if so, update the seats counts on those flights and release
		the mutual exclusive access and let another goroutine be able to access those flights
	*/
	for i := 0; i < len(flights) && bookable; i++ {
		if flights[i].SeatsLeft < seatsTobook {
			bookable = false
		}
	}

	for i := 0; i < len(flights) && bookable; i++ {
		flights[i].SeatsLeft -= seatsTobook
	}

	for _, f := range flights {
		f.Locker.Unlock()
	}

	return bookable
}

func NewFlight(origin, dest string) *Flight {
	return &Flight{
		Origin:    origin,
		Dest:      dest,
		SeatsLeft: 200,
		Locker:    NewSpinLock(),
	}
}
