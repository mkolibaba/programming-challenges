package commands

import (
	"fmt"
	"github.com/mkolibaba/programming-challenges/sudoku-solver-go/puzzle"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func NewApproxCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "approx",
		Short: "Runs approx",
		Long:  "Runs super duper approx",
		Args:  cobra.ExactArgs(1),
		Run:   approxHandler,
	}
}

func approxHandler(cmd *cobra.Command, args []string) {
	n, _ := strconv.Atoi(args[0])

	//logToFile, _ := tea.LogToFile("approx.txt", "debug")
	//defer logToFile.Close()

	//log.Printf("Solving %d sudokus\n", n)
	sudokus := make([]*puzzle.Sudoku, n)
	for i := 0; i < n; i++ {
		sudoku := puzzle.NewFromKaggle(i)
		sudoku.Solve()
		//log.Printf("Sudoku %d solved: moves = %d, time = %s\n", i, sudoku.Moves, sudoku.TimeElapsed)
		sudokus[i] = sudoku
	}

	var timeWasted time.Duration
	var averageMoves int
	for _, sudoku := range sudokus {
		timeWasted += sudoku.TimeElapsed
		averageMoves += sudoku.Moves
	}
	averageMoves = averageMoves / n

	fmt.Printf("Total time: %s\n", timeWasted)
	fmt.Printf("Average moves: %d\n", averageMoves)
	fmt.Printf("Average time: %s\n", timeWasted/time.Duration(n))
}
