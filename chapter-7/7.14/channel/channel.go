package channel

import (
	"container/list"
	"sync"

	"7.14/semaphore"
)

type Channel[M any] struct {
	/*
	   using 2 semaphores to block for separate senarios.
	   the capacitySemaphore will block when our buffer is filled,
	   while the sizeSemaphore will block when our buffer is empty
	   and there is no messages to consume.
	*/
	capacitySemaphore *semaphore.Semaphore

	sizeSemaphore *semaphore.Semaphore

	mutex sync.Mutex

	buffer *list.List // standard library queue implementation
}

func NewChannel[M any](capacity int) *Channel[M] {
	return &Channel[M]{
		capacitySemaphore: semaphore.NewSemaphore(capacity),
		sizeSemaphore:     semaphore.NewSemaphore(0),
		buffer:            list.New(),
	}
}

func (c *Channel[M]) Send(message M) {
	/*
	   when attempting to send a message through the channel, we
	   attempt to acquire a permit on the capacitySemaphore.
	   there's always an initial capacity permit based on the capacity
	   value passed in during channel instantiation.

	   if our buffer is filled up, then the capacitySemaphore's permit is 0.
	   any goroutine that attempts to send more messages through this channel
	   will be blocked until other goroutines calls .Release() to increment
	   the permits in the capacitySemaphore to allow .Send() to acquire permits.
	*/
	c.capacitySemaphore.Acquire()

	c.mutex.Lock()
	c.buffer.PushBack(message)
	c.mutex.Unlock()

	/*
	   increment the sizeSemaphore permits whenever we're able to add a message to the buffer.
	   this will allow goroutines that attempts to consume messages be able to acquire a permit
	   and consume a message.
	   if there is no message in the buffer, then the initial permit counter for sizeSemaphore is 0
	   so any goroutines that calls .Receive() then they are blocked until the permit is incremented
	   here.
	*/
	c.sizeSemaphore.Release()
}

func (c *Channel[M]) Receive() M {
	/*
	   increments the capacitySemaphore's permit count by 1 since we're consuming a message
	   off the buffer. this will allow another goroutine to send 1 more message to the buffer.
	*/
	c.capacitySemaphore.Release()

	/*
	   goroutines that call .Receive() here will block if there is no message in the buffer
	   as determined by the sizeSemaphore permit count. this will be 0 initially until
	   another goroutine calls .Send() to send a message and increment the sizeSemaphore's permit
	   counter by 1.
	*/
	c.sizeSemaphore.Acquire()

	c.mutex.Lock()
	v := c.buffer.Remove(c.buffer.Front()).(M)
	c.mutex.Unlock()

	return v
}
