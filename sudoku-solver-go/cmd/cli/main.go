package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"github.com/mkolibaba/programming-challenges/sudoku-solver-go/render"
	"os"
)

func main() {
	p := tea.NewProgram(render.NewModel("1.txt"))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

/**
sudoku resources:
https://qqwing.com/?spm=a2ty_o01.29997173.0.0.5426c921YjptHN
*/
