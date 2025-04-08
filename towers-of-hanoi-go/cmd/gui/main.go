package main

import (
	"fyne.io/fyne/v2"
	"image/color"
	"math"
)

var (
	// palette: https://www.color-hex.com/color-palette/1057778
	colorRed    = createColor(224, 125, 69)
	colorGreen  = createColor(143, 149, 98)
	colorBlack  = createColor(0, 0, 0)
	colorYellow = createColor(218, 152, 60)
	colorWhite  = createColor(244, 221, 184)

	rodsXPositioning = []float32{0.2, 0.5, 0.8}

	diskTopY    = platformY - float32(math.Ceil(float64(platformThickness)/2)) - diskHeight*3
	diskBottomY = platformY - float32(math.Ceil(float64(platformThickness)/2)) - diskHeight

	windowSize = fyne.NewSize(windowWidth, windowHeight)
)

const (
	windowWidth  = float32(600)
	windowHeight = float32(400)

	gap = float32(50)

	platformThickness = float32(15)
	platformX1        = gap
	platformX2        = windowWidth - gap
	platformY         = float32(200)
	platformLength    = platformX2 - platformX1

	rodY1        = platformY - diskHeight*4
	rodY2        = platformY
	rodThickness = platformThickness * 2 / 3

	diskHeight = float32(40)

	smallDiskWidth  = float32(80)
	mediumDiskWidth = float32(100)
	largeDiskWidth  = float32(120)
)

func main() {
	gui := NewGUI()
	gui.Run()
}

func createColor(r, g, b uint8) color.NRGBA {
	return color.NRGBA{R: r, G: g, B: b, A: 255}
}
