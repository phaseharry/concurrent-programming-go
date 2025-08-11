package main

import "sync"

type Arbitrator struct {
	accountsInUse map[string]bool
	cond          *sync.Cond
}

func NewArbitrator() *Arbitrator {
	return &Arbitrator{
		accountsInUse: map[string]bool{},
		cond:          sync.NewCond(&sync.Mutex{}),
	}
}

/*
acquiring the lock on the Arbitrator so only one goroutine
can lock multiple accounts at a time
*/
func (a *Arbitrator) LockAccounts(ids ...string) {
	a.cond.L.Lock()
	/*
	   have the goroutine that is attempting to lock multiple accounts
	   continuing trying until it's able to get a lock on all of the required
	   accounts it needs
	*/
	for allAvailable := false; !allAvailable; {
		allAvailable = true
		for _, id := range ids {
			/*
				if true, then current requested accountId is in-use by another goroutine, so
				give up the exclusive access to the arbitrator lock and try again later.
				only mark requested accountIds as in use, if all the requested accountIds
				are available
			*/
			if a.accountsInUse[id] {
				allAvailable = false
				a.cond.Wait()
				break
			}
		}
	}
	/*
	   if all the requestedIds are available, marks them as in-use, making it unavailable for other
	   goroutines to make changes to those accounts while current goroutine has exclusive access to all of them
	*/
	for _, id := range ids {
		a.accountsInUse[id] = true
	}
	a.cond.L.Unlock()
}

func (a *Arbitrator) UnlockAccounts(ids ...string) {
	a.cond.L.Lock()
	/*
		marking all ids as not in-use and then broadcasting & then unlocking
		so any blocked goroutine in the LockAccounts method can attempt to
		acquire exclusive access on its requested accountIds
	*/
	for _, id := range ids {
		a.accountsInUse[id] = false
	}
	a.cond.Broadcast()
	a.cond.L.Unlock()
}
