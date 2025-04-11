package main

import (
	"fmt"
	"github.com/mkolibaba/programming-challenges/sudoku-solver-go/puzzle"
	"os"
	"time"
)

func main() {
	fmt.Print("Enter number of sudokus: ")

	var n int
	fmt.Fscanf(os.Stdin, "%d", &n)

	sudokus := make([]*puzzle.Sudoku, n)
	for i := 0; i < n; i++ {
		sudoku := puzzle.NewRandomFromKaggle()
		sudoku.Solve()
		sudokus[i] = sudoku
	}

	var averageTimeDuration time.Duration
	var averageMoves int
	for _, sudoku := range sudokus {
		averageTimeDuration += sudoku.TimeElapsed
		averageMoves += sudoku.Moves
	}
	averageTimeDuration = averageTimeDuration / time.Duration(n)
	averageMoves = averageMoves / n

	fmt.Printf("Average moves: %d\n", averageMoves)
	fmt.Printf("Average time: %s\n", averageTimeDuration)
}
