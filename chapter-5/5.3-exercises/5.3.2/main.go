package main

import (
	"fmt"
	"sync"
	"time"
)

/*
stops waiting for players to connect if all did not connect to the
game within 10 seconds and start the game.
*/
func main() {
	cond := sync.NewCond(&sync.Mutex{})
	cancel := false
	go timeout(cond, &cancel)
	/*
	   stating that the players needed to join as 5 so we will
	   always start the game through timeout and not because all
	   players joined.
	*/
	playersInGame := 5

	// only connecting 4 players
	for i := range 4 {
		go playerHandler(cond, &playersInGame, i, &cancel)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(60 * time.Second)
}

/*
waiting 10 seconds and then setting
the cancel value to true and then broadcast it.
*/
func timeout(cond *sync.Cond, cancel *bool) {
	time.Sleep(10 * time.Second)
	cond.L.Lock()
	*cancel = true
	cond.Broadcast()
	cond.L.Unlock()
}

func playerHandler(cond *sync.Cond, playersRemaining *int, playerId int, cancel *bool) {
	cond.L.Lock()
	fmt.Println(playerId, ": Connected")

	// decrementing a player count because player connected.
	*playersRemaining--

	/*
	   will never be called since we're artifically not connecting all players.
	   all player goroutine will be blocked until the timeout happens and
	   set cancel = true and unlock all player goroutines.
	*/
	if *playersRemaining == 0 {
		cond.Broadcast()
	}

	/*
	   when all the players goroutines get a Broadcast() event,
	   it will check the playersRemaining value. If it is greater
	   than 0 then that means we're still waiting for players to connect.
	   player goroutines also stop waiting if gets a Broadcast() event
	   and cancel is true due to the timeout.
	*/
	for *playersRemaining > 0 && !*cancel {
		fmt.Println(playerId, ": Waiting for more players")
		cond.Wait()
	}
	cond.L.Unlock()
	if *cancel {
		fmt.Println(playerId, ": Game cancelled")
	} else {
		fmt.Println("All players connected. Ready player", playerId)
	}
}
