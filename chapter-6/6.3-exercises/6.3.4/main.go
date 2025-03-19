package main

import (
	"fmt"
	"math/rand"
	"sync"
)

const matrixSize = 3

func main() {
	var matrixA, matrixB, result [matrixSize][matrixSize]int
	generateRandMatrix(&matrixA)
	generateRandMatrix(&matrixB)

	fmt.Println("matrices before multiplication")
	printMatrices(&matrixA, &matrixB, &result)

	wg := sync.WaitGroup{}
	wg.Add(matrixSize)

	for row := range matrixSize {
		go rowMultiply(&matrixA, &matrixB, &result, row, &wg)
	}
	wg.Wait()
	fmt.Println("matrices after multiplication")
	printMatrices(&matrixA, &matrixB, &result)
}

func generateRandMatrix(matrix *[matrixSize][matrixSize]int) {
	for row := range matrixSize {
		for col := range matrixSize {
			matrix[row][col] = rand.Intn(10) - 5
		}
	}
}

func rowMultiply(
	matrixA, matrixB, result *[matrixSize][matrixSize]int, row int, wg *sync.WaitGroup,
) {
	for col := range matrixSize {
		sum := 0
		for i := range matrixSize {
			sum += matrixA[row][i] * matrixB[i][col]
		}
		result[row][col] = sum
	}
	wg.Done()
}

func printMatrices(matrixA, matrixB, result *[matrixSize][matrixSize]int) {
	for row := range matrixSize {
		fmt.Println(matrixA[row], matrixB[row], result[row])
	}
}
