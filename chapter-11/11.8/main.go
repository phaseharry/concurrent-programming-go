package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	accounts := []BankAccount{
		*NewBankAccount("Sam"),
		*NewBankAccount("Paul"),
		*NewBankAccount("Amy"),
		*NewBankAccount("Mia"),
	}
	arb := NewArbitrator()

	total := len(accounts)
	for i := range 4 {
		wg.Add(1)
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
				accounts[from].Transfer(&accounts[to], 10, eId, arb)
			}
			fmt.Println(eId, "COMPLETE")
			wg.Done()
		}(i)
	}
	wg.Wait()
}
