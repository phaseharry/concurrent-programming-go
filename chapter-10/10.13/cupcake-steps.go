package listing10_13

import (
	"fmt"
	"time"
)

const (
	OVEN_TIME            = 5
	EVERYTHING_ELSE_TIME = 2
)

// helper functions used to compare pipelining to sequential processing
func PrepareTray(trayNumber int) string {
	fmt.Println("Preparing empty tray", trayNumber)
	time.Sleep(EVERYTHING_ELSE_TIME * time.Second)
	return fmt.Sprintf("tray number %d", trayNumber)
}

func Mixture(tray string) string {
	fmt.Println("Pouring cupcake Mixture in", tray)
	time.Sleep(EVERYTHING_ELSE_TIME * time.Second)
	return fmt.Sprintf("cupcake in %s", tray)
}

func Bake(mixture string) string {
	fmt.Println("Baking", mixture)
	time.Sleep(OVEN_TIME * time.Second)
	return fmt.Sprintf("baked %s", mixture)
}

func AddToppings(bakedCupCake string) string {
	fmt.Println("Adding topping to", bakedCupCake)
	time.Sleep(EVERYTHING_ELSE_TIME * time.Second)
	return fmt.Sprintf("topping on %s", bakedCupCake)
}

func Box(finishedCupCake string) string {
	fmt.Println("Boxing", finishedCupCake)
	time.Sleep(EVERYTHING_ELSE_TIME * time.Second)
	return fmt.Sprintf("%s boxed", finishedCupCake)
}
