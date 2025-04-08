package hanoi

import (
	"errors"
	"slices"
)

type Disk int

const (
	Small Disk = iota
	Medium
	Large

	WelcomeInfo = `Welcome to the game Tower of Hanoi!
The rules of the Tower of Hanoi are:
1. You should move all the disks to the last rod.
2. Only one disk can be moved at a time.
3. A disk can only be placed on a larger disk or an empty rod.
Game can be solved in 2^N - 1 moves for N disks.`
)

var (
	ErrMoveSameRod       = errors.New("moving disk to the same rod is senseless")
	ErrMoveNoDisks       = errors.New("there is no disks on rod")
	ErrMoveToSmallerDisk = errors.New("moving disk to smaller disk is forbidden")
)

type Pile []Disk

type Game struct {
	Piles [3]Pile
	Moves int
}

func (g *Game) Move(from, to int) error {
	if from == to {
		return ErrMoveSameRod
	}

	fromPile := g.Piles[from]
	pileTo := g.Piles[to]
	if len(fromPile) == 0 {
		return ErrMoveNoDisks
	}

	disk := fromPile[0]
	if len(pileTo) != 0 && pileTo[0] < disk {
		return ErrMoveToSmallerDisk
	}

	g.Piles[from] = fromPile[1:]
	g.Piles[to] = append(Pile{disk}, pileTo...)
	g.Moves++
	return nil
}

func (g *Game) IsFinished() bool {
	return len(g.Piles[0]) == 0 &&
		len(g.Piles[1]) == 0 &&
		len(g.Piles[2]) == 3 &&
		slices.Equal(g.Piles[2], Pile{Small, Medium, Large})
}

func NewGame() *Game {
	return &Game{
		Piles: [3]Pile{
			{Small, Medium, Large},
			nil,
			nil,
		},
	}
}
