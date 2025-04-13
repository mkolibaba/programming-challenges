package commands

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"github.com/mkolibaba/programming-challenges/sudoku-solver-go/puzzle"
	"github.com/mkolibaba/programming-challenges/sudoku-solver-go/render"
	"github.com/spf13/cobra"
	"os"
)

func NewSolverCmd() *cobra.Command {
	solverCmd := &cobra.Command{
		Use:   "solver",
		Short: "Runs interactive sudoku solver",
	}
	kaggleCmd := &cobra.Command{
		Use:   "kaggle",
		Short: "Solves sudoku from file kaggle-sudoku.csv",
		Run:   solverKaggleHandler,
	}
	kaggleCmd.Flags().IntP("number", "n", -1, "Puzzle number")
	solverCmd.AddCommand(kaggleCmd)
	fileCmd := &cobra.Command{
		Use:   "compact-view",
		Short: "Solves sudoku from compact-view file",
		Run:   solverCompactViewHandler,
	}
	fileCmd.Flags().StringP("filename", "f", "", "File name")
	fileCmd.MarkFlagRequired("filename")
	solverCmd.AddCommand(fileCmd)
	return solverCmd
}

func solverKaggleHandler(cmd *cobra.Command, args []string) {
	n, _ := cmd.Flags().GetInt("number")
	var sudoku *puzzle.Sudoku
	if n != -1 {
		sudoku = puzzle.NewFromKaggle(n)
	} else {
		sudoku = puzzle.NewRandomFromKaggle()
	}
	initTui(sudoku)
}

func solverCompactViewHandler(cmd *cobra.Command, args []string) {
	filename, _ := cmd.Flags().GetString("filename")
	initTui(puzzle.NewFromFile(filename))
}

func initTui(sudoku *puzzle.Sudoku) {
	p := tea.NewProgram(render.NewModel(sudoku))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

/**
sudoku resources:
https://qqwing.com/?spm=a2ty_o01.29997173.0.0.5426c921YjptHN
*/
