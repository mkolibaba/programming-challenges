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
)

type GUI struct {
	game             *hanoi.Game
	application      fyne.App
	welcomeInfoLabel *widget.Label
	movesLabel       *widget.Label
	infoLabel        *widget.Label
	rods             []fyne.CanvasObject
	disks            []*Disk
}

func NewGUI() *GUI {
	// инициализируем application в самом начале, перед созданием компонентов
	application := app.New()
	game := hanoi.NewGame()
	movesLabel := widget.NewLabel(fmt.Sprintf("Moves: %d", game.Moves))
	infoLabel := widget.NewLabel("")
	rods := createRods()
	disks := createDisks(game)

	gui := &GUI{
		game:             game,
		application:      application,
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
	application := gui.application
	window := application.NewWindow("Towers of Hanoi")

	header := gui.initializeHeaderContainer()
	drawings := gui.initializeDrawingsContainer()

	content := container.NewVBox(header, gui.infoLabel, drawings)

	window.SetContent(content)

	window.Resize(windowSize)
	window.CenterOnScreen()
	window.ShowAndRun()
}

func (gui *GUI) finish() {
	gui.infoLabel.SetText("Game is finished. You won!")
	for i, _ := range gui.disks {
		gui.disks[i].Release()
	}
}

func (gui *GUI) afterMove(err error) {
	gui.movesLabel.SetText(fmt.Sprintf("Moves: %d", gui.game.Moves))
	if err != nil {
		gui.infoLabel.SetText(fmt.Sprintf("Error: %v", err))
	} else if gui.game.IsFinished() {
		gui.finish()
	} else {
		gui.infoLabel.SetText("")
	}
}

func createRods() (rods []fyne.CanvasObject) {
	for i := 0; i < 3; i++ {
		rods = append(rods, NewRod(i))
	}
	return
}

func createDisks(game *hanoi.Game) (disks []*Disk) {
	for pileIdx, pile := range game.Piles {
		for diskIdx, disk := range pile {
			rect := NewDisk(
				diskBottomY-diskHeight*float32(len(pile)-diskIdx-1),
				disk,
				pileIdx,
			)
			disks = append(disks, rect)
		}
	}
	return
}

// TODO: переписать метод
func setOnDragEvent(disk *Disk, gui *GUI) {
	disk.OnDragEnd = func() (fyne.Position, bool) {
		return gui.moveDisk(disk)
	}
}

func (gui *GUI) moveDisk(disk *Disk) (fyne.Position, bool) {
	i, rod := disk.FindRodByCurrentPosition(gui.rods)
	if i != -1 {
		err := gui.game.Move(disk.OnRod, i)
		gui.afterMove(err)
		if err != nil {
			return fyne.Position{}, false
		}
		x := rod.Position().X - disk.Size().Width/2
		y := diskBottomY - diskHeight*float32(len(gui.game.Piles[i])-1)
		disk.OnRod = i
		return fyne.NewPos(x, y), true
	}

	return fyne.Position{}, false
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
