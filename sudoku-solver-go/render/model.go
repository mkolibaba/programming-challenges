package render

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mkolibaba/programming-challenges/sudoku-solver-go/puzzle"
	"strings"
	"time"
)

type Model struct {
	sudoku    puzzle.Sudoku
	sub       chan struct{}
	stopwatch stopwatch.Model
	keyMap    KeyMap
	help      help.Model
}

type KeyMap struct {
	start key.Binding
	quit  key.Binding
}

func NewModel(filename string) Model {
	return Model{
		sudoku:    puzzle.NewFromFile(filename),
		stopwatch: stopwatch.NewWithInterval(time.Nanosecond),
		sub:       make(chan struct{}),
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
			return m, tea.Sequence(
				m.stopwatch.Start(),
				createSolveCmd(m),
				createListenCmd(m),
			)
		case key.Matches(msg, m.keyMap.quit):
			return m, tea.Quit
		}
	case SudokuSolvedMsg:
		return m, m.stopwatch.Stop()
	}
	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	builder := &strings.Builder{}
	builder.WriteString(m.sudoku.String())

	fmt.Fprintf(builder, " Elapsed: %s\n", m.stopwatch.View())
	builder.WriteString(m.renderHelp())
	builder.WriteString("\n")
	return builder.String()
}

func (m Model) renderHelp() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keyMap.start,
		m.keyMap.quit,
	})
}

type SudokuSolvedMsg struct{}

func createSolveCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		go func() {
			m.sudoku.Solve()
			m.sub <- struct{}{}
		}()
		return nil
	}
}

func createListenCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		return SudokuSolvedMsg(<-m.sub)
	}
}
