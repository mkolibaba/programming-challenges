package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"image/color"
)

const (
	rodY1        = platformY - diskHeight*4
	rodY2        = platformY
	rodThickness = platformThickness * 2 / 3
)

var (
	rodsXPositioning = []float32{0.2, 0.5, 0.8}
)

type Rod canvas.Line

func NewRod(n int) *canvas.Line {
	rod := canvas.NewLine(color.Black)
	rod.StrokeWidth = rodThickness
	x := CalculateRodX(n)
	rod.Position1 = fyne.NewPos(x, rodY1)
	rod.Position2 = fyne.NewPos(x, rodY2)
	return rod
}

// TODO: по идее можно найти лучшее место для этого
func CalculateRodX(n int) float32 {
	return platformLength*rodsXPositioning[n] + platformX1
}
