package channel

import (
	"container/list"
	"sync"
)

// channel implemented with Condition Variable
type Channel[M any] struct {
	cond     sync.Cond
	buffer   *list.List // standard library queue implementation
	capacity int
}

func NewChannel[M any](capacity int) *Channel[M] {
	return &Channel[M]{
		cond:     *sync.NewCond(&sync.Mutex{}),
		buffer:   list.New(),
		capacity: capacity,
	}
}

func (c *Channel[M]) Send(message M) {
	c.cond.L.Lock()

	for c.buffer.Len() == c.capacity {
		c.cond.Wait()
	}

	c.buffer.PushBack(message)
	c.cond.Broadcast()
	c.cond.L.Unlock()
}

func (c *Channel[M]) Receive() M {
	c.cond.L.Lock()

	/*
		ensure we can receive messages through an unbufferred channel
		by increasing capacity every time .Receive() is called. (ex. capacity == 0)
		so when .Send() is called, it will be blocked until there's a goroutine
		that calls .Receive() and increment that capacity to 1
	*/
	c.capacity++

	/*
	   broadcasting here so any blocked .Send() calls can attempt to get lock
	   and push message to buffer.
	*/
	c.cond.Broadcast()

	for c.buffer.Len() == 0 {
		c.cond.Wait()
	}

	c.capacity--
	v := c.buffer.Remove(c.buffer.Front()).(M)
	c.cond.L.Unlock()

	return v
}
