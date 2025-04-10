package render

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mkolibaba/programming-challenges/sudoku-solver-go/puzzle"
	"strings"
	"sync"
	"time"
)

type Model struct {
	sudoku    *puzzle.Sudoku
	stopwatch stopwatch.Model
	keyMap    KeyMap
	help      help.Model
	done      bool
}

type KeyMap struct {
	start key.Binding
	quit  key.Binding
}

func NewModel(sudoku *puzzle.Sudoku) Model {
	return Model{
		sudoku:    sudoku,
		stopwatch: stopwatch.NewWithInterval(time.Microsecond),
		keyMap: KeyMap{
			start: key.NewBinding(
				key.WithKeys("x"),
				key.WithHelp("x", "start"),
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
				m.stopwatch.Start(),
				createSolveCmd(m, &wg),
				createListenCmd(&wg),
			)
		case key.Matches(msg, m.keyMap.quit):
			return m, tea.Quit
		}
	case SudokuSolvedMsg:
		m.done = true
		return m, m.stopwatch.Stop()
	case stopwatch.TickMsg:
		var cmd tea.Cmd
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	builder := &strings.Builder{}
	builder.WriteString(m.sudoku.String())

	if m.done {
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
