package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	"os"
	"strings"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	sudoku    [][]int
	sub       chan struct{}
	spinner   spinner.Model
	responses int
	quitting  bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			return m, tea.Batch(
				m.spinner.Tick,
				listenForActivity(m), // generate activity
				waitForActivity(m),   // wait for activity
			)
		case "q":
			m.quitting = true
			return m, tea.Quit
		}
	case responseMsg:
		m.responses++                // record external activity
		return m, waitForActivity(m) // wait for next event
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

//func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	switch msg := msg.(type) {
//	case tea.KeyMsg:
//		switch msg.String() {
//		case "q":
//			return m, tea.Quit
//		}
//	}
//	return m, nil
//}

func (m model) View() string {

	builder := &strings.Builder{}
	builder.WriteString(prettyPrint(m.sudoku))

	builder.WriteString("\n")

	fmt.Fprintf(builder, "\n %s Events received: %d\n\n Press any key to exit\n", m.spinner.View(), m.responses)
	// TODO: bubble help
	if m.quitting {
		builder.WriteString("\n")
	}
	return builder.String()
}

func initialModel() model {
	return model{
		sudoku:  readFromFile("input/1.txt"),
		spinner: spinner.New(),
		sub:     make(chan struct{}),
	}
}

func solve(sudoku [][]int) {
	i, j := getNextBlank(sudoku)
	if i == -1 && j == -1 {
		return
	}
	solveNext(sudoku, i, j)
}

func solveNext(sudoku [][]int, row, column int) bool {
	for possibleValue := 1; possibleValue <= 9; possibleValue++ {
		sudoku[row][column] = possibleValue
		if checkCell(sudoku, row, column) {
			i, j := getNextBlank(sudoku)
			if i == -1 && j == -1 {
				return true
			}
			if solveNext(sudoku, i, j) {
				return true
			}
		}
	}
	// ни одно из значений не подходит
	sudoku[row][column] = 0
	return false
}

//func solveOld(sudoku [][]int, i, j int) {
//	for row, _ := range sudoku {
//		if row > 1 {
//			continue
//		}
//
//		for column, _ := range sudoku[row] {
//			if !isBlank(sudoku, row, column) {
//				continue
//			}
//			for possibleValue := 1; possibleValue <= 9; possibleValue++ {
//				sudoku[row][column] = possibleValue
//				if checkCell(sudoku, row, column) {
//					break
//				}
//			}
//		}
//	}
//}

func getNextBlank(sudoku [][]int) (int, int) {
	for i, row := range sudoku {
		for j, cell := range row {
			if cell == 0 {
				return i, j
			}
		}
	}
	return -1, -1
}

func isBlank(sudoku [][]int, i, j int) bool {
	return sudoku[i][j] == 0
}

/**
sudoku resources:
https://qqwing.com/?spm=a2ty_o01.29997173.0.0.5426c921YjptHN
*/
