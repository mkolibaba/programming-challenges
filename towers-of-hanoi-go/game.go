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
	moves     int
	moveError error // TODO: как-то сомнительно
}

func (g *Game) Move(from, to int) {
	if from == to {
		g.moveError = fmt.Errorf("moving disk to the same rod is senseless")
		return
	}

	fromPile := g.Piles[from]
	if len(fromPile) == 0 {
		g.moveError = fmt.Errorf("there is no disks on rod %d", from)
		return
	}

	// TODO: error processing
	disk := g.Piles[from][0]
	g.Piles[from] = g.Piles[from][1:]
	g.Piles[to] = append(Pile{disk}, g.Piles[to]...)
	g.moves++
	g.moveError = nil
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
