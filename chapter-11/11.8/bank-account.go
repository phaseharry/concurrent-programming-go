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

func (src *BankAccount) Transfer(to *BankAccount, amount int, exId int, arb *Arbitrator) {
	fmt.Printf("%d locking %s and %s\n", exId, src.id, to.id)
	arb.LockAccounts(src.id, to.id)
	src.balance -= amount
	to.balance += amount
	arb.UnlockAccounts(src.id, to.id)
	fmt.Printf("%d unlocked %s and %s\n", exId, src.id, to.id)
}
