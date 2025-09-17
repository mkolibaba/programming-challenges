package main

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

func main() {
	listPath := flag.String("l", "", "Path of list files to download")
	flag.Parse()

	if *listPath == "" {
		fmt.Println("You must specify a path to a list of files to download")
		return
	}

	downloads, err := FromList(*listPath)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	m := initialModel(downloads)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("error during program execution: %v", err)
	}
}

func initialModel(downloads []*Download) Model {
	return Model{
		spinner:   spinner.New(spinner.WithSpinner(spinner.Globe)),
		Downloads: downloads,
	}
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.spinner.Tick,
		doTick(100 * time.Millisecond),
	}
	for _, d := range m.Downloads {
		cmds = append(cmds, d.Run)
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case ErrMsg:
		m.err = msg
		return m, nil

	case TickMsg:
		return m, doTick(100 * time.Millisecond)

	case DoneMsg:
		for _, d := range m.Downloads {
			if d.Status != Done {
				return m, nil
			}
		}
		m.done = true
		return m, tea.Quit

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m Model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	var lines []string
	for _, d := range m.Downloads {
		var s string
		switch d.Status {
		case Fetching:
			s = fmt.Sprintf("%s Fetching %s", m.spinner.View(), d.URL)
		case InProgress:
			s = fmt.Sprintf(
				InProgressTemplate,
				m.spinner.View(), d.Name, d.DurationHumanized(), d.SpeedHumanized(), d.Downloaded, d.Size,
			)
		case Done:
			s = fmt.Sprintf(DoneTemplate, d.Name, d.DurationHumanized(), d.SpeedHumanized(), d.Size)
		}
		lines = append(lines, s)
	}

	if m.done {
		lines = append(lines, "🐣 All downloads completed successfully!\n")
	}
	if m.quitting {
		lines = append(lines, "\n")
	}

	str := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return str
}

type TickMsg time.Time

func doTick(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type Model struct {
	spinner   spinner.Model
	quitting  bool
	err       error
	Downloads []*Download
	done      bool
}

type ErrMsg error

type DoneMsg struct{}

var InProgressTemplate = `%s Downloading %s
  ⌚ Duration: %s
  🐢 Speed: %s
  💾 Total: %s / %s`

var DoneTemplate = `✅ Downloaded %s
  ⌚ Duration: %s
  🐢 Speed: %s
  💾 Total: %s`
