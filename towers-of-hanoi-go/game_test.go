package main_test

import (
	"github.com/mkolibaba/programming-challenges/towers-of-hanoi-go"
	"testing"
)

func TestGame(t *testing.T) {
	t.Run("test game initial state", func(t *testing.T) {
		game := main.NewGame()

		if game == nil {
			t.Errorf("game is nil")
		}

		if game.IsFinished() {
			t.Errorf("game is finished in its initial state")
		}
	})
}
