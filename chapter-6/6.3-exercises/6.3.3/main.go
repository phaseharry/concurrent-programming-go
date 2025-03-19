package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const matrixSize = 1000

func main() {
	var matrixA, matrixB, result [matrixSize][matrixSize]int
	generateRandMatrix(&matrixA)
	generateRandMatrix(&matrixB)

	sequentialStart := time.Now()
	sequentialMatrixMultiply(&matrixA, &matrixB, &result)
	sequentialDuration := time.Since(sequentialStart)

	wg := sync.WaitGroup{}
	wg.Add(matrixSize)
	concurrentStart := time.Now()
	for row := range matrixSize {
		go rowMultiply(&matrixA, &matrixB, &result, row, &wg)
	}
	wg.Wait()
	concurrentDuration := time.Since(concurrentStart)

	fmt.Printf(
		"Sequential matrix multiplication took %v while concurrent matrix multiplication took %v.",
		sequentialDuration,
		concurrentDuration,
	)
}

func generateRandMatrix(matrix *[matrixSize][matrixSize]int) {
	for row := range matrixSize {
		for col := range matrixSize {
			// generate a random number between -5 and 4 for each
			// element within the matrix
			matrix[row][col] = rand.Intn(10) - 5
		}
	}
}

func sequentialMatrixMultiply(matrixA, matrixB, result *[matrixSize][matrixSize]int) {
	for row := range matrixSize {
		for col := range matrixSize {
			sum := 0
			for i := range matrixSize {
				sum += matrixA[row][i] * matrixB[i][col]
			}
			result[row][col] = sum
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
