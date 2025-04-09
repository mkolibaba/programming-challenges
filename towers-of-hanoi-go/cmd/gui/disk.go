package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/mkolibaba/programming-challenges/towers-of-hanoi-go"
	"image/color"
)

// TODO: связать диски с Game
type Disk struct {
	widget.BaseWidget
	rect          *canvas.Rectangle
	boundPosition fyne.Position
	OnRod         int
	blocked       bool // TODO: нельзя двигать, если не верхний диск
	OnDragEnd     func() (fyne.Position, bool)
}

func NewDisk(y float32, diskType hanoi.Disk, onRod int) *Disk {
	var clr color.Color
	var w float32
	switch diskType {
	case hanoi.Small:
		clr, w = colorRed, smallDiskWidth
	case hanoi.Medium:
		clr, w = colorGreen, mediumDiskWidth
	case hanoi.Large:
		clr, w = colorYellow, largeDiskWidth
	}

	x := CalculateRodX(onRod) - w/2

	rectangle := canvas.NewRectangle(clr)
	disk := &Disk{
		rect:          rectangle,
		boundPosition: fyne.NewPos(x, y),
		blocked:       false,
		OnRod:         onRod,
	}
	disk.ExtendBaseWidget(disk)
	disk.Resize(fyne.NewSize(w, diskHeight))
	disk.Move(fyne.NewPos(x, y))
	return disk
}

func (d *Disk) Dragged(event *fyne.DragEvent) {
	if !d.blocked {
		d.Move(fyne.NewPos(d.Position().X+event.Dragged.DX, d.Position().Y+event.Dragged.DY))
	}
}

func (d *Disk) DragEnd() {
	if d.OnDragEnd != nil && !d.blocked {
		newPosition, ok := d.OnDragEnd()
		if ok {
			// фиксируем новое местоположение
			d.Move(newPosition)
			d.boundPosition = d.Position()
		} else {
			// возвращаем на прежнее место
			d.Move(d.boundPosition)
		}
	}
}

func (d *Disk) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(d.rect)
}

// FindRodByCurrentPosition находит rod в зависимости от своего текущего положения
func (d *Disk) FindRodByCurrentPosition(rods []fyne.CanvasObject) (int, fyne.CanvasObject) {
	getBounds := func(obj fyne.CanvasObject) (left, right, top, bottom float32) {
		return obj.Position().X, obj.Position().X + obj.Size().Width, obj.Position().Y, obj.Position().Y + obj.Size().Height
	}
	diskLeft, diskRight, diskTop, diskBottom := getBounds(d)

	for i, rod := range rods {
		rodLeft, rodRight, rodTop, rodBottom := getBounds(rod)
		if rodRight > diskLeft && rodLeft < diskRight && rodBottom > diskTop && rodTop < diskBottom {
			return i, rod
		}
	}

	return -1, nil
}

func (d *Disk) Release() {
	d.blocked = true
}
