package render

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mkolibaba/programming-challenges/sudoku-solver-go/activity"
	"github.com/mkolibaba/programming-challenges/sudoku-solver-go/puzzle"
	"strings"
)

type Model struct {
	Sudoku    [][]int
	Sub       chan struct{}
	Spinner   spinner.Model
	responses int
	Quitting  bool
}

func NewModel(path string) Model {
	return Model{
		Sudoku:  puzzle.ReadFromFile(path),
		Spinner: spinner.New(),
		Sub:     make(chan struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			return m, tea.Batch(
				m.Spinner.Tick,
				activity.ListenForActivity(m.Sub), // generate activity
				activity.WaitForActivity(m.Sub),   // wait for activity
			)
		case "q":
			m.Quitting = true
			return m, tea.Quit
		}
	case activity.ResponseMsg:
		m.responses++                             // record external activity
		return m, activity.WaitForActivity(m.Sub) // wait for next event
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
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

func (m Model) View() string {

	builder := &strings.Builder{}
	builder.WriteString(puzzle.PrettyPrint(m.Sudoku))

	builder.WriteString("\n")

	fmt.Fprintf(builder, "\n %s Events received: %d\n\n Press any key to exit\n", m.Spinner.View(), m.responses)
	// TODO: bubble help
	if m.Quitting {
		builder.WriteString("\n")
	}
	return builder.String()
}
