package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"github.com/mkolibaba/programming-challenges/sudoku-solver-go/puzzle"
	"github.com/mkolibaba/programming-challenges/sudoku-solver-go/render"
	"math/rand"
	"os"
)

func main() {
	//text := puzzle.NewFromFile("main.txt")
	kaggle := puzzle.NewFromKaggle(rand.Intn(1_000_000))
	p := tea.NewProgram(render.NewModel(kaggle))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

/**
sudoku resources:
https://qqwing.com/?spm=a2ty_o01.29997173.0.0.5426c921YjptHN
*/
