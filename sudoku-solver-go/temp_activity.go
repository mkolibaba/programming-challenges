package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"math/rand"
	"time"
)

// A message used to indicate that activity has occurred. In the real world (for
// example, chat) this would contain actual data.
type responseMsg struct{}

// Simulate a process that sends events at an irregular interval in real time.
// In this case, we'll send events on the channel at a random interval between
// 100 to 1000 milliseconds. As a command, Bubble Tea will run this
// asynchronously.
func listenForActivity(m model) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Millisecond * time.Duration(rand.Int63n(100)+100)) // nolint:gosec
			m.sub <- struct{}{}
		}
	}
}

// A command that waits for the activity on a channel.
func waitForActivity(m model) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-m.sub)
	}
}
