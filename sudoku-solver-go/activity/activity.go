package activity

import (
	tea "github.com/charmbracelet/bubbletea"
	"math/rand"
	"time"
)

// A message used to indicate that activity has occurred. In the real world (for
// example, chat) this would contain actual data.
type ResponseMsg struct{}

// Simulate a process that sends events at an irregular interval in real time.
// In this case, we'll send events on the channel at a random interval between
// 100 to 1000 milliseconds. As a command, Bubble Tea will run this
// asynchronously.
func ListenForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Millisecond * time.Duration(rand.Int63n(100)+100)) // nolint:gosec
			sub <- struct{}{}
		}
	}
}

// A command that waits for the activity on a channel.
func WaitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return ResponseMsg(<-sub)
	}
}
