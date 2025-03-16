package main

import (
	"fmt"
	"math/rand"

	"6.16/barrier"
)

const matrixSize = 3

func main() {
	var matrixA, matrixB, result [matrixSize][matrixSize]int

	/*
	   setting barrier size to be matrixSize + 1
	*/
	barr := barrier.NewBarrier(matrixSize + 1)

	/*
	   starting the matrix multiplication. Even if all 3 rows have called rowMultiply,
	   they will be initially blocked since our barrier waitSize is matrixSize + 1 (4).
	   the 3 rowMultiply goroutines will be blocked until the main goroutine generates
	   random matrices for A and B and calls barrier.Wait() as well.
	*/
	for row := range matrixSize {
		go rowMultiply(&matrixA, &matrixB, &result, row, barr)
	}

	for range 4 {
		generateRandMatrix(&matrixA)
		generateRandMatrix(&matrixB)
		fmt.Println("matrices before multiplication")
		printMatrices(&matrixA, &matrixB, &result)
		/*
			blocks main goroutine until all rowMultiply goroutines
			and main goroutine has called .Wait()
			once all called, all the rowMultiply goroutines can process
			and calculate the matrix multiplication results for their
			assigned rows.
		*/
		barr.Wait()
		/*
		   block again to wait for all rowMultiply goroutines to finish
		   their processing and call .Wait(). Once that happens waitSize will
		   reach the barrier size (4) and unblock all goroutines and reset the waitSize to 0
		*/
		barr.Wait()
		fmt.Println("matrices after multiplication")
		printMatrices(&matrixA, &matrixB, &result)
	}
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

/*
function meant to be used as a goroutine so all rows can calculate their
matrix mutliplication results concurrently.
*/
func rowMultiply(
	matrixA, matrixB, result *[matrixSize][matrixSize]int,
	row int,
	barr *barrier.Barrier,
) {
	/*
		infinite loop to continuous generate matrix multiplications of rand matrices
		for each row (0, 1, 2). Calls .Wait() at the start of each iteration to ensure
		that all existing goroutines have reached this point and called .Wait(). Once
		all goroutines have then allow them all to concurrently process and calculate
		a matrix multiplication result for the current row
	*/
	for {
		barr.Wait()
		for col := range matrixSize {
			sum := 0
			for i := range matrixSize {
				sum += matrixA[row][i] * matrixB[i][col]
			}
			result[row][col] = sum
		}
		/*
		   have the barrier call .Wait() for each rowMultiply goroutine so all goroutines
		   will be blocked until all rowMultiply() goroutines have finished calculation for their row's
		   respective matrix multiplications before moving on
		*/
		barr.Wait()
	}
}

// printing input matrices and the resulting matrix
func printMatrices(matrixA, matrixB, result *[matrixSize][matrixSize]int) {
	for row := range matrixSize {
		fmt.Println(matrixA[row], matrixB[row], result[row])
	}
}
