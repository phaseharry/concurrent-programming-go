package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	player1, player2, player3, player4 := player(), player(), player(), player()
	playersRemaining := 4

	for player1 != nil || player2 != nil || player3 != nil || player4 != nil {
		select {
		case p1, p1Active := <-player1:
			if !p1Active {
				player1 = nil
				playersRemaining -= 1
				fmt.Println("Player 1 left the game, Remaining players: ", playersRemaining)
			} else {
				fmt.Println("Player 1:", p1)
			}
		case p2, p2Active := <-player2:
			if !p2Active {
				player2 = nil
				playersRemaining -= 1
				fmt.Println("Player 2 left the game, Remaining players: ", playersRemaining)
			} else {
				fmt.Println("Player 2:", p2)
			}
		case p3, p3Active := <-player3:
			if !p3Active {
				player3 = nil
				playersRemaining -= 1
				fmt.Println("Player 3 left the game, Remaining players: ", playersRemaining)
			} else {
				fmt.Println("Player 3:", p3)
			}
		case p4, p4Active := <-player4:
			if !p4Active {
				player4 = nil
				playersRemaining -= 1
				fmt.Println("Player 4 left the game, Remaining players: ", playersRemaining)
			} else {
				fmt.Println("Player 4:", p4)
			}

		}
	}
}

func player() chan string {
	output := make(chan string)
	count := rand.Intn(100)
	move := []string{"UP", "DOWN", "LEFT", "RIGHT"}
	go func() {
		defer close(output)
		for range count {
			output <- move[rand.Intn(4)]
			d := time.Duration(rand.Intn(200))
			time.Sleep(d * time.Millisecond)
		}
	}()
	return output
}
