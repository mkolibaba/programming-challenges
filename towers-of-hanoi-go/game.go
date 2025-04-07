package main

import (
	"fmt"
	"slices"
)

type Disk int

const (
	Small Disk = iota
	Medium
	Large
)

type Pile []Disk

type Game struct {
	Piles     [3]Pile
	Moves     int
	MoveError error // TODO: как-то сомнительно
}

func (g *Game) Move(from, to int) {
	if from == to {
		g.MoveError = fmt.Errorf("moving disk to the same rod is senseless")
		return
	}

	fromPile := g.Piles[from]
	pileTo := g.Piles[to]
	if len(fromPile) == 0 {
		g.MoveError = fmt.Errorf("there is no disks on rod %d", from)
		return
	}

	// TODO: error processing
	disk := fromPile[0]
	if len(pileTo) != 0 && pileTo[0] < disk {
		g.MoveError = fmt.Errorf("moving disk to smaller disk is forbidden")
		return
	}

	g.Piles[from] = fromPile[1:]
	g.Piles[to] = append(Pile{disk}, pileTo...)
	g.Moves++
	g.MoveError = nil
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
