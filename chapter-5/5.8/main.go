package main

import (
	"fmt"
	"sync"
	"time"
)

/*
Example demonstrating .Broadcast() signaling multiple goroutines that they can stop waiting and
attempt to acquire the lock continue boardcasting.

In this example, it simulates waiting for 4 players (goroutines) to connect to a game server (main goroutine).
We spawn 4 child goroutines are the first 3 is blocked because we only call .Boardcast() and stop all waiting
when the "playersInGames" counter hit 0 indicating that all 4 players are ready. Once that happens the the game "starts".

We don't use .Signal() here because there are multiple goroutines we want to unblock and have continue processing.
If we did use .Signal() instead of Broadcast(), then 2 of the child goroutines / players would be blocked,
leading to a deadlock since there won't be another goroutine spawned to call the exact amount of .Signal()
calls to unblock their waits. This will lead to deadlocks.
*/
func main() {
	cond := sync.NewCond(&sync.Mutex{})
	playersInGame := 4
	for playerId := range 4 {
		go playerHandler(cond, &playersInGame, playerId)
		time.Sleep(1 * time.Second)
	}
}

/*
4 playerHandler goroutines will be spawned with playerIds 1, 2, 3, 4 in any possible order.
Will order it as 1, 2, 3, 4 to make it easier to understands.
1.

	player1 calls, acquires lock and decrements playersInGame from 4 -> 3.


	It seems that the we're still waiting for players (playersInGame is not zero) so it calls .Wait().


	This will release the lock and have this goroutine wait until there's a Signal() or Broadcast() call.

2.

	player2 calls, acquires lock and decrements playersInGame from 3 -> 2.
	same as above.

3.

	player3 calls, acquires lock and decrements playersInGame from 2 -> 1.
	same as above.

4.

	player4 calls, acquires lock and decrements playersInGame from 1 -> 0.
	Since playersInGame is 0, this goroutine calls .Boardcast(), letting all goroutines
	that are currently waiting that they can attempt to acquire lock when it's available.
	Current goroutine will unlock and let other goroutines attempt to acquire it.

5.

	goroutine for players 1, 2, and 3 will attempt to acquire lock & run through the look again.
	Since playersInGame == 0, they will execute past the for loop and able to unlock and
	finish the goroutine. All players are connected and ready to start game.
*/
func playerHandler(cond *sync.Cond, playersInGame *int, playerId int) {
	cond.L.Lock()
	fmt.Println(playerId, ": Connected")
	*playersInGame -= 1
	if *playersInGame == 0 {
		cond.Broadcast()
	}

	for *playersInGame > 0 {
		fmt.Println(playerId, ": Waiting for more players")
		cond.Wait()
	}
	cond.L.Unlock()
	fmt.Println("All players connected. Ready player", playerId)
}
