package main

import (
	"fmt"
	"sync"
)

type BankAccount struct {
	id      string
	balance int
	mutex   sync.Mutex
}

func NewBankAccount(id string) *BankAccount {
	return &BankAccount{
		id:      id,
		balance: 100,
		mutex:   sync.Mutex{},
	}
}

func (src *BankAccount) Transfer(to *BankAccount, amount int, exId int) {
	fmt.Printf("%d locking %s's account\n", exId, src.id)
	src.mutex.Lock()
	fmt.Printf("%d locking %s's account\n", exId, src.id)
	to.mutex.Lock()

	src.balance -= amount
	to.balance += amount

	src.mutex.Unlock()
	to.mutex.Unlock()
	fmt.Printf("%d unlocked %s and %s\n", exId, src.id, to.id)
}
