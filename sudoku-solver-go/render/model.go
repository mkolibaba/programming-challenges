package render

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mkolibaba/programming-challenges/sudoku-solver-go/puzzle"
	"strings"
	"sync"
)

type Model struct {
	sudoku *puzzle.Sudoku
	keyMap KeyMap
	help   help.Model
}

type KeyMap struct {
	start   key.Binding
	restart key.Binding
	quit    key.Binding
}

func NewModel(sudoku *puzzle.Sudoku) Model {
	return Model{
		sudoku: sudoku,
		keyMap: KeyMap{
			start: key.NewBinding(
				key.WithKeys("x"),
				key.WithHelp("x", "start"),
			),
			restart: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "restart"),
				key.WithDisabled(),
			),
			quit: key.NewBinding(
				key.WithKeys("q"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.start):
			var wg sync.WaitGroup
			wg.Add(1)
			m.keyMap.start.SetEnabled(false)
			return m, tea.Batch(
				createSolveCmd(m, &wg),
				createListenCmd(&wg),
			)
		case key.Matches(msg, m.keyMap.restart):
			m.keyMap.restart.SetEnabled(false)
			m.keyMap.start.SetEnabled(true)
			m.sudoku = puzzle.NewRandomFromKaggle()
		case key.Matches(msg, m.keyMap.quit):
			return m, tea.Quit
		}
	case SudokuSolvedMsg:
		m.keyMap.restart.SetEnabled(true)
	}
	return m, nil
}

func (m Model) View() string {
	builder := &strings.Builder{}
	builder.WriteString(m.sudoku.String())

	if m.sudoku.IsSolved() {
		builder.WriteString("\n Done!")
		fmt.Fprintf(builder, "\n Time elapsed: %s\n", m.sudoku.TimeElapsed)
	}
	builder.WriteString(m.renderHelp())
	builder.WriteString("\n")
	return builder.String()
}

func (m Model) renderHelp() string {
	return "\n " + m.help.ShortHelpView([]key.Binding{
		m.keyMap.start,
		m.keyMap.restart,
		m.keyMap.quit,
	})
}

type SudokuSolvedMsg struct{}

func createSolveCmd(m Model, c *sync.WaitGroup) tea.Cmd {
	return func() tea.Msg {
		go func() {
			defer c.Done()
			m.sudoku.Solve()
		}()
		return nil
	}
}

func createListenCmd(c *sync.WaitGroup) tea.Cmd {
	return func() tea.Msg {
		c.Wait()
		return SudokuSolvedMsg{}
	}
}
