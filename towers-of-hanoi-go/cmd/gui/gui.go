package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mkolibaba/programming-challenges/towers-of-hanoi-go"
	"image/color"
)

type GUI struct {
	game             *hanoi.Game
	welcomeInfoLabel *widget.Label
	movesLabel       *widget.Label
	infoLabel        *widget.Label
	rods             []fyne.CanvasObject
	disks            []*Disk
}

func NewGUI() *GUI {
	game := hanoi.NewGame()
	movesLabel := widget.NewLabel(fmt.Sprintf("Moves: %d", game.Moves))
	infoLabel := widget.NewLabel("")
	rods := createRods()
	disks := createDisks(game, rods) // TODO: get rid of passing rods (maybe)

	gui := &GUI{
		game:             game,
		welcomeInfoLabel: widget.NewLabel(hanoi.WelcomeInfo),
		movesLabel:       movesLabel,
		infoLabel:        infoLabel,
		rods:             rods,
		disks:            disks,
	}
	// TODO: корректно поправить
	for i, _ := range disks {
		setOnDragEvent(disks[i], gui)
	}
	return gui
}

func (gui *GUI) Run() {
	application := app.New()
	window := application.NewWindow("Towers of Hanoi")

	header := gui.initializeHeaderContainer()
	drawings := gui.initializeDrawingsContainer()

	content := container.NewVBox(header, gui.infoLabel, drawings)

	window.SetContent(content)

	window.Resize(windowSize)
	window.CenterOnScreen()
	window.ShowAndRun()
}

func (gui *GUI) MoveDiskOld(err error) {
	gui.movesLabel.SetText(fmt.Sprintf("Moves: %d", gui.game.Moves))
	if err != nil {
		gui.infoLabel.SetText(fmt.Sprintf("Error: %v", err))
	} else if gui.game.IsFinished() {
		gui.infoLabel.SetText("Game is finished. You won!")
		for i, _ := range gui.disks {
			gui.disks[i].Release()
		}
	} else {
		gui.infoLabel.SetText("")
	}
}

func createRods() []fyne.CanvasObject {
	rods := make([]fyne.CanvasObject, 3)
	for i := 0; i < 3; i++ {
		rod := canvas.NewLine(color.Black)
		rod.StrokeWidth = rodThickness
		x := platformLength*rodsXPositioning[i] + platformX1
		rod.Position1 = fyne.NewPos(x, rodY1)
		rod.Position2 = fyne.NewPos(x, rodY2)
		rods[i] = rod
	}
	return rods
}

func createDisks(game *hanoi.Game, rods []fyne.CanvasObject) []*Disk {
	var disks []*Disk

	for pileIdx, pile := range game.Piles {
		for diskIdx, disk := range pile {
			var clr color.Color
			var w float32
			switch disk {
			case hanoi.Small:
				clr, w = colorRed, smallDiskWidth
			case hanoi.Medium:
				clr, w = colorGreen, mediumDiskWidth
			case hanoi.Large:
				clr, w = colorYellow, largeDiskWidth
			}
			rect := NewDisk(
				clr,
				w,
				diskHeight,
				rods[pileIdx].Position().X-w/2,
				diskBottomY-diskHeight*float32(len(pile)-diskIdx-1),
			)
			disks = append(disks, rect)
		}
	}

	return disks
}

// TODO: переписать метод
func setOnDragEvent(disk *Disk, gui *GUI) {
	disk.OnDragEnd = func() bool {
		return gui.moveDisk(disk)
	}
}

func (gui *GUI) moveDisk(disk *Disk) bool {
	for i, rod := range gui.rods {
		if disk.IntersectsWith(rod) {
			err := gui.game.Move(disk.GetCurrentRod(gui.rods), i)
			gui.MoveDiskOld(err)
			if err != nil {
				return false
			}
			x := gui.rods[i].Position().X - disk.Size().Width/2
			y := diskBottomY - diskHeight*float32(len(gui.game.Piles[i])-1)
			disk.Move(fyne.NewPos(x, y))
			return true
		}
	}
	return false
}

func (gui *GUI) initializeHeaderContainer() *fyne.Container {
	return container.NewHBox(gui.welcomeInfoLabel, layout.NewSpacer(), gui.movesLabel)
}

func (gui *GUI) initializeDrawingsContainer() *fyne.Container {
	// платформа рисуется здесь, т.к. с ней нет никакого взаимодействия
	platform := canvas.NewLine(colorBlack)
	platform.StrokeWidth = platformThickness
	platform.Position1 = fyne.NewPos(platformX1, platformY)
	platform.Position2 = fyne.NewPos(platformX2, platformY)

	return container.NewWithoutLayout(
		platform,
		gui.rods[0], gui.rods[1], gui.rods[2],
		gui.disks[0], gui.disks[1], gui.disks[2],
	)
}
