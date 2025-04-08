package hanoi_test

import (
	"errors"
	"github.com/mkolibaba/programming-challenges/towers-of-hanoi-go"
	"testing"
)

func TestGame(t *testing.T) {
	t.Run("game initial state", func(t *testing.T) {
		game := hanoi.NewGame()

		if game == nil {
			t.Errorf("game is nil")
		}

		if game.IsFinished() {
			t.Errorf("game is finished in its initial state")
		}
	})
	t.Run("moving disk to the same rod gives error", func(t *testing.T) {
		game := hanoi.NewGame()

		err := game.Move(0, 0)

		if err == nil {
			t.Errorf("should return error when moving disk to the same rod")
		}
		if !errors.Is(err, hanoi.ErrMoveSameRod) {
			t.Errorf("should return error of type ErrMoveSameRod")
		}
	})
	t.Run("moving disk from empty rod", func(t *testing.T) {
		game := hanoi.NewGame()

		err := game.Move(1, 0)

		if err == nil {
			t.Errorf("should return error when moving disk from empty rod")
		}
		if !errors.Is(err, hanoi.ErrMoveNoDisks) {
			t.Errorf("should return error of type ErrMoveNoDisks")
		}
	})
	t.Run("moving disk to smaller disks gives error", func(t *testing.T) {
		game := hanoi.NewGame()

		game.Move(0, 1)
		err := game.Move(0, 1)

		if err == nil {
			t.Errorf("should return error when moving disk to smaller disks")
		}
		if !errors.Is(err, hanoi.ErrMoveToSmallerDisk) {
			t.Errorf("should return error of type ErrMoveToSmallerDisk")
		}
	})
	t.Run("should properly track moves", func(t *testing.T) {
		game := hanoi.NewGame()

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
		game := hanoi.NewGame()

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
	t.Run("should recover from exception", func(t *testing.T) {
		game := hanoi.NewGame()

		err := game.Move(1, 0)

		if err == nil {
			t.Errorf("should return error")
		}

		err = game.Move(0, 1)

		if err != nil {
			t.Errorf("should recover from error, got %v", err)
		}
		if game.Moves != 1 {
			t.Errorf("should increment moves, got %d", game.Moves)
		}
	})
}
