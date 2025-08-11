package main

import (
	"fmt"
	"math/rand"
	"time"

	listing11_3_4 "github.com/phaseharry/concurrent-programming-go/chapter-11/11.4"
)

/*
demoing the possibility of a deadlock happening in which 2 goroutines compete for the same exclusive resource.
ex.
- goroutine 1 attempts to process a transaction from Sam -> Paul
- goroutine 2 attempts to process a transaction from Paul -> Amy
- goroutine 3 attempts to process a transaction from Amy -> Mia
- goroutine 4 attempts to process a transaction from Mia -> Sam
if all 3 goroutines operate concurrently and
1. acquires lock to Sam
2. acquires lock to Paul
3. acquires lock to Amy
4. acquires lock to Mia
then
1. will be blocked since it cannot acquire a lock to Paul afterwards since 2 has it
2. will be blocked since it cannot acquire a lock to Amy since 3 has it
3. will be blocked since it cannot acquire a lock to Mia since 4 has it
4. will be blocked since it cannot acquire a lock to Sam since 1 has it
*/
func main() {
	accounts := []listing11_3_4.BankAccount{
		*listing11_3_4.NewBankAccount("Sam"),
		*listing11_3_4.NewBankAccount("Paul"),
		*listing11_3_4.NewBankAccount("Amy"),
		*listing11_3_4.NewBankAccount("Mia"),
	}

	total := len(accounts)
	for i := range 4 {
		// create 4 goroutines simulating multiple transactions happening in parallel (4)
		go func(eId int) {
			// for each goroutine, execute 1000 transactions
			for j := 1; j < 1000; j++ {
				// generate random from & to accounts to simulate a transaction
				from, to := rand.Intn(total), rand.Intn(total)
				/*
					while the from & to generated values are the same, keep generating a new to value until it's a different
					account id
				*/
				for from == to {
					to = rand.Intn(total)
				}
				// transfering the account
				accounts[from].Transfer(&accounts[to], 10, eId)
			}
			fmt.Println(eId, "COMPLETE")
		}(i)
	}
	/*
	   let program run for 60 seconds or until a deadlock is detected from the 4 goroutines
	   running concurrently, whichever happens first
	*/
	time.Sleep(60 * time.Second)
}
