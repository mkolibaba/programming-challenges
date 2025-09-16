package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
	"time"
)

const link = "https://github.com/portapps/postman-portable/releases/download/11.62.7-64/postman-portable-win64-11.62.7-64-setup.exe"
const link2 = "https://github.com/pocketbase/pocketbase/releases/download/v0.30.0/pocketbase_0.30.0_windows_amd64.zip"

func main() {
	logger := zap.Must(zap.NewDevelopment()).Sugar()

	m := initialModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		logger.Fatalf("error during program execution: %v", err)
	}
}

func initialModel() Model {
	return Model{
		spinner:   spinner.New(spinner.WithSpinner(spinner.Globe)),
		Downloads: []*Download{NewDownload(link), NewDownload(link2)},
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
		case InProgress:
			s = fmt.Sprintf(
				InProgressTemplate,
				m.spinner.View(), d.Name, d.DurationHumanized(), d.SpeedHumanized(), d.DownloadedHumanized(), d.SizeHumanized(),
			)
		case Done:
			s = fmt.Sprintf(
				DoneTemplate,
				d.Name, d.DurationHumanized(), d.SpeedHumanized(), d.SizeHumanized(),
			)
		}
		lines = append(lines, s)
	}

	str := lipgloss.JoinVertical(lipgloss.Left, lines...)
	if m.quitting {
		return str + "\n"
	}
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
}

type ErrMsg error

var InProgressTemplate = `%s Downloading %s
  ⌚ Duration: %s
  🕑 Speed: %s
  💾 Total: %s / %s`

var DoneTemplate = `✅ Downloaded %s
  ⌚ Duration: %s
  🕑 Speed: %s
  💾 Total: %s`
