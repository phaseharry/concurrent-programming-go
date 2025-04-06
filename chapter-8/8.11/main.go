package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
demonstrating setting channels to nil to have cases
that is blocked on receiving events from those channels
to not be used after it was set to nil
*/
func main() {
	sales := generateAmounts(50)
	expenses := generateAmounts(40)
	endOfDayAmount := 0
	for sales != nil || expenses != nil {
		select {
		case sale, saleChannelActive := <-sales:
			/*
			   using the 2nd parameter of receiving messages from channel. this is a boolean value
			   indicating whether or not a channel has been closed or not. if it is false, then we
			   know we finished sending messages through this queue, so we will set the sales channel
			   to nil so the case

			   case, moreData := sales

			   will not listened to anymore / will not block and only
			   expense case will process if it has not been set to nil as well
			*/
			if saleChannelActive {
				fmt.Println("Sale of:", sale)
				endOfDayAmount += sale
			} else {
				sales = nil
			}
		case expense, expenseChannelActive := <-expenses:
			if expenseChannelActive {
				fmt.Println("Expense of:", expense)
				endOfDayAmount -= expense
			} else {
				expenses = nil
			}
		}
	}
	fmt.Println("End of day profit and loss:", endOfDayAmount)
}

/*
returns a channel and starts a goroutine that sends n messages
to that channel and then closes it.
*/
func generateAmounts(n int) <-chan int {
	amounts := make(chan int)
	go func() {
		defer close(amounts)
		for range n {
			amounts <- rand.Intn(100) + 1
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return amounts
}
