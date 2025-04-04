package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	t, _ := strconv.Atoi(os.Args[1]) // reads timeout value from environment when app is run
	messages := sendMsgAfter(3 * time.Second)
	timeoutDuration := time.Duration(t) * time.Second
	fmt.Printf("Waiting for message for %d seconds...\n", t)

	/*
	   demonstrating timing out a blocked consumption of a message from channel.
	   if the case where we consume a message from the "messages" channel does not
	   have a message within the specified timeoutDuration, then we stop being blocked
	   on that case and process the timeout case.
	   time.After creates a channel that accepts time types. Once timeoutDuration is up,
	   it sends a time message through that channel so our timeout case runs
	*/
	select {
	case msg := <-messages:
		fmt.Println("Message received:", msg)
	case tNow := <-time.After(timeoutDuration):
		fmt.Println("Timed out. Waited until:", tNow.Format("15:04:05"))
	}

}

/*
create a receiver channel and returns it for another
function to consume messages from.
also creates a separate goroutine that sends messages
to that queue every x seconds
*/
func sendMsgAfter(seconds time.Duration) <-chan string {
	messages := make(chan string)

	go func() {
		time.Sleep(seconds)
		messages <- "Hello"
	}()

	return messages
}
