package main_test

import (
	"github.com/mkolibaba/programming-challenges/towers-of-hanoi-go"
	"testing"
)

func TestGame(t *testing.T) {
	t.Run("game initial state", func(t *testing.T) {
		game := main.NewGame()

		if game == nil {
			t.Errorf("game is nil")
		}

		if game.IsFinished() {
			t.Errorf("game is finished in its initial state")
		}
	})
	t.Run("moving disk to the same rod gives error", func(t *testing.T) {
		game := main.NewGame()

		game.Move(0, 0)

		if game.MoveError == nil {
			t.Errorf("should return error when moving disk to the same rod")
		}
	})
	t.Run("moving disk from empty rod", func(t *testing.T) {
		game := main.NewGame()

		game.Move(1, 0)

		if game.MoveError == nil {
			t.Errorf("should return error when moving disk from empty rod")
		}
	})
	t.Run("moving disk to smaller disks gives error", func(t *testing.T) {
		game := main.NewGame()

		game.Move(0, 1)
		game.Move(0, 1)

		if game.MoveError == nil {
			t.Errorf("should return error when moving disk to smaller disks")
		}
	})
	t.Run("should properly track moves", func(t *testing.T) {
		game := main.NewGame()

		game.Move(0, 2)
		game.Move(2, 1)
		game.Move(1, 2)
		game.Move(2, 1)
		game.Move(1, 2)

		got := game.Moves
		want := 5

		if got != want {
			t.Errorf("got %d moves, want %d", got, want)
		}
	})
	t.Run("should win game", func(t *testing.T) {
		game := main.NewGame()

		game.Move(0, 2)
		game.Move(0, 1)
		game.Move(2, 1)
		game.Move(0, 2)
		game.Move(1, 0)
		game.Move(1, 2)
		game.Move(0, 2)

		if !game.IsFinished() {
			t.Errorf("game should be finished")
		}
	})
}
