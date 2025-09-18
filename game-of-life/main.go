package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

var initial = [][]bool{
	{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	{false, false, true, false, false, false, false, false, false, false, false, false, true, false, false},
	{false, false, false, true, false, false, false, false, false, false, false, true, false, false, false},
	{false, false, true, true, false, false, false, false, false, false, false, true, true, false, false},
	{false, false, false, false, false, false, true, false, false, false, false, false, false, false, false},
	{false, false, false, false, false, true, true, false, false, true, true, false, false, false, false},
	{false, false, false, false, false, true, false, false, false, false, true, false, false, false, false},
	{false, true, false, false, false, false, false, false, false, false, false, false, false, false, false},
	{false, false, false, false, false, true, false, false, false, false, true, false, false, false, false},
	{false, false, false, false, false, true, true, false, false, true, true, false, false, false, false},
	{true, false, false, false, false, false, true, false, false, false, false, false, false, false, false},
	{false, true, true, true, false, false, false, false, false, false, false, true, true, false, false},
	{false, true, false, true, false, false, false, false, false, false, false, true, false, false, false},
	{false, false, true, false, false, false, false, false, false, false, false, false, true, false, false},
	{false, false, false, false, false, false, true, false, false, false, false, false, false, false, false},
}

var glider15x15 = [][]bool{
	{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, true, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false, true, false, false, false, false, false, false},
	{false, false, false, false, false, false, true, true, true, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
}

func main() {
	m := Model{
		state: glider15x15,
		iter:  1,
	}
	program := tea.NewProgram(m)
	if _, err := program.Run(); err != nil {
		fmt.Println(err)
	}
}

type Model struct {
	state [][]bool
	iter  int
}

func (m Model) Init() tea.Cmd {
	return doTick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			return m, nil
		}

	case TickMsg:
		m.state = newGen(m.state)
		m.iter++
		return m, doTick()

	default:
		return m, nil
	}
}

func (m Model) View() string {
	str := fmt.Sprintf("Iter %d:\n", m.iter)
	for _, r := range m.state {
		for _, b := range r {
			str += boolToString(b) + " "
		}
		str += "\n"
	}

	return str
}

type TickMsg struct{}

func doTick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

func newGen(s [][]bool) [][]bool {
	n := make([][]bool, len(s))
	for i := range s {
		n[i] = make([]bool, len(s[i]))
		for j := range s[i] {
			n[i][j] = isAlive(i, j, s)
		}
	}
	return n
}

func makeACopy(s [][]bool) [][]bool {
	n := make([][]bool, len(s))
	for i := range s {
		n[i] = make([]bool, len(s[i]))
		for j := range s[i] {
			n[i][j] = s[i][j]
		}
	}
	return n
}

func isAlive(i, j int, c [][]bool) bool {
	aliveNeighbours := 0

	for row := i - 1; row <= i+1; row++ {
		for col := j - 1; col <= j+1; col++ {
			if row == i && col == j {
				continue
			}
			if row >= 0 && len(c) > row && col >= 0 && len(c[row]) > col && c[row][col] {
				aliveNeighbours++
			}
		}
	}

	if c[i][j] {
		return aliveNeighbours == 2 || aliveNeighbours == 3
	}
	return aliveNeighbours == 3
}

func boolToString(b bool) string {
	if b {
		return "X"
	}
	return "."
}
