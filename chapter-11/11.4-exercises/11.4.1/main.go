package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
)

type Player struct {
	name  string
	score int
	mutex sync.Mutex
}

func incrementScores(players []*Player, increment int) {
	/*
		sorting players by their unique name values and then acquiring lock
		on individual players. this ensures that all goroutines will always acquire
		exclusive access to players in the same order so no 2 goroutines can acquire partial access
		to resources but be blocked from acquiring other exclusive resources that are held
		by other goroutines
	*/
	sort.Slice(players, func(a, b int) bool {
		return players[a].name < players[b].name
	})
	for _, player := range players {
		player.mutex.Lock()
	}
	for _, player := range players {
		player.score += increment
	}
	for _, player := range players {
		player.mutex.Unlock()
	}
}

func main() {
	players := []*Player{
		{"harry", 0, sync.Mutex{}},
		{"jason", 0, sync.Mutex{}},
		{"arnold", 0, sync.Mutex{}},
		{"clark", 0, sync.Mutex{}},
	}

	wg := sync.WaitGroup{}

	for range 1000 {
		n := rand.Intn(len(players)) + 1

		rand.Shuffle(
			len(players),
			func(i, j int) { players[i], players[j] = players[j], players[i] },
		)

		wg.Add(1)
		sublist := make([]*Player, n)
		copy(sublist, players[:n])

		go func(players []*Player) {
			incrementScores(players, 10)
			wg.Done()
		}(sublist)
	}

	wg.Wait()
	for _, player := range players {
		fmt.Printf("Score has %s is %d\n", player.name, player.score)
	}
}
