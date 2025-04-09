package main

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/programming-challenges/towers-of-hanoi-go"
	"os"
	"slices"
	"strings"
)

var (
	// palette: https://www.color-hex.com/color-palette/1057814
	styleRed    = lipgloss.NewStyle().Foreground(lipgloss.Color("#dc6f4e"))
	styleGreen  = lipgloss.NewStyle().Foreground(lipgloss.Color("#bce36e"))
	styleYellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#e2d052"))

	styleSelect = lipgloss.NewStyle().Foreground(lipgloss.Color("#d000ff"))
	styleError  = lipgloss.NewStyle().Foreground(lipgloss.Color("#c00000"))

	renderedDisks = map[hanoi.Disk]string{
		hanoi.Small:  "  " + styleRed.Render("ooo") + "  ",
		hanoi.Medium: " " + styleGreen.Render("ooooo") + " ",
		hanoi.Large:  styleYellow.Render("ooooooo"),
	}
	renderedRod = "   |   "

	platformBase     = "xxxx%sxxxxxxx%sxxxxxxx%sxxxx"
	cursorTemplate   = strings.ReplaceAll(platformBase, "x", " ")
	renderedPlatform = strings.ReplaceAll(platformBase, "%s", "x")

	boolToFloat64 = map[bool]float64{true: 1, false: 0}
)

type model struct {
	game        *hanoi.Game
	cursor      int
	from        int
	lastError   error
	completeBar progress.Model
}

func (m *model) moveLeft() {
	m.cursor = (m.cursor + len(m.game.Piles) - 1) % len(m.game.Piles)
	if m.cursor == m.from {
		m.moveLeft()
	}
}

func (m *model) moveRight() {
	m.cursor = (m.cursor + 1) % len(m.game.Piles)
	if m.cursor == m.from {
		m.moveRight()
	}
}

func (m *model) moveDisk() {
	if m.from != -1 {
		err := m.game.Move(m.from, m.cursor)
		if err != nil {
			m.lastError = err
		}
		m.from = -1
	} else {
		m.from = m.cursor
		m.moveRight()
		m.lastError = nil
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	gameActions := map[string]func(){
		"left":  m.moveLeft,
		"right": m.moveRight,
		"enter": m.moveDisk,
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "right", "enter":
			if !m.game.IsFinished() {
				gameActions[msg.String()]()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	return NewTUIBuilder().
		Line(hanoi.WelcomeInfo).
		Formatted("\nMoves: %d\n", m.game.Moves).
		Line(renderCompleteBar(m)).
		EmptyLine().
		Line(renderRodsAndPiles(m)).
		Line(renderCursor(m)).
		With(func(b *TUIBuilder) {
			if m.lastError != nil {
				b.Formatted("%s %s\n", styleError.Render("Invalid move:"), m.lastError)
			}
		}).
		With(func(b *TUIBuilder) {
			if m.game.IsFinished() {
				b.Line("Game is finished. You won!")
				b.Line("Press q to quit.")
			}
		}).
		Build()
}

func renderCursor(m model) string {
	cursors := make([]any, 3)
	for i := 0; i < len(cursors); i++ {
		cursors[i] = " "
	}
	cursors[m.cursor] = "^"
	if m.from != -1 {
		cursors[m.from] = styleSelect.Render("^")
	}

	return fmt.Sprintf(cursorTemplate, cursors...)
}

func renderCompleteBar(m model) string {
	pile := slices.Clone(m.game.Piles[2])
	slices.Reverse(pile)
	plen := len(pile)
	correctDisksCount := boolToFloat64[plen > 0 && pile[0] == hanoi.Large] +
		boolToFloat64[plen > 1 && pile[1] == hanoi.Medium] +
		boolToFloat64[plen > 2 && pile[2] == hanoi.Small]
	bar := m.completeBar.ViewAs(correctDisksCount / 3)
	return fmt.Sprintf("Complete: %s", bar)
}

func renderRodsAndPiles(m model) string {
	buf := &bytes.Buffer{}
	game := m.game
	height := 4
	for i := 0; i < height; i++ {
		// first row always empty
		var rows []string
		for _, pile := range game.Piles {
			mappedIndex := i - (height - len(pile))
			if mappedIndex >= 0 {
				rows = append(rows, renderedDisks[pile[mappedIndex]])
			} else {
				rows = append(rows, renderedRod)
			}
		}
		fmt.Fprintf(buf, " %s \n", strings.Join(rows, " "))
	}
	fmt.Fprint(buf, renderedPlatform)
	return buf.String()
}

func initialModel() model {
	bar := progress.New(progress.WithScaledGradient("#6b19b3", "#d000ff"), progress.WithWidth(20))
	return model{
		game:        hanoi.NewGame(),
		from:        -1,
		completeBar: bar,
	}
}
