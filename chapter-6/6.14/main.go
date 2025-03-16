package main

import (
	"fmt"
	"math/rand"
)

const matrixSize = 3

func main() {
	var matrixA, matrixB, result [matrixSize][matrixSize]int
	generateRandMatrix(&matrixA)
	generateRandMatrix(&matrixB)
	fmt.Println("matrices before matrix multiply")
	printMatrices(&matrixA, &matrixB, &result)
	matrixMultiply(&matrixA, &matrixB, &result)
	fmt.Println("matrices after matrix multiply")
	printMatrices(&matrixA, &matrixB, &result)
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

// implementation of matrix mutliplication between 2 matrixs that are both n*n size
func matrixMultiply(matrixA, matrixB, result *[matrixSize][matrixSize]int) {
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

// printing input matrices and the resulting matrix
func printMatrices(matrixA, matrixB, result *[matrixSize][matrixSize]int) {
	for row := range matrixSize {
		fmt.Println(matrixA[row], matrixB[row], result[row])
	}
}
