package main

import (
	"fmt"

	listing10_13 "github.com/phaseharry/concurrent-programming-go/chapter-10/10.13"
)

func main() {
	/*
	   simulating sequential processing of 1 person baking 10
	   trays of cupcakes.

	   run: time go run cupcake-sequential.go

	   to get total time needed to make 10 batches
	*/
	for i := range 10 {
		result := listing10_13.Box(
			listing10_13.AddToppings(
				listing10_13.Bake(
					listing10_13.Mixture(
						listing10_13.PrepareTray(i),
					),
				),
			),
		)
		fmt.Println("Acceping", result)
	}
}
