package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

type Rod canvas.Line

// TODO: связать диски с Game
type Disk struct {
	widget.BaseWidget
	rect          *canvas.Rectangle
	boundPosition fyne.Position
	blocked       bool        // TODO: нельзя двигать, если не верхний диск
	OnDragEnd     func() bool // TODO: возможно, должна возвращать координаты
}

func NewDisk(color color.Color, width, height, x, y float32) *Disk {
	rectangle := canvas.NewRectangle(color)
	disk := &Disk{
		rect:          rectangle,
		boundPosition: fyne.NewPos(x, y),
		blocked:       false,
	}
	disk.ExtendBaseWidget(disk)
	disk.Resize(fyne.NewSize(width, height))
	disk.Move(fyne.NewPos(x, y)) // TODO: есть warning, что сначала нужно стартануть приложение, а потом двигать
	return disk
}

func (d *Disk) Dragged(event *fyne.DragEvent) {
	if !d.blocked {
		d.Move(fyne.NewPos(d.Position().X+event.Dragged.DX, d.Position().Y+event.Dragged.DY))
	}
}

func (d *Disk) DragEnd() {
	if d.OnDragEnd != nil && !d.blocked {
		if d.OnDragEnd() {
			// фиксируем новое местоположение
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

func (d *Disk) IntersectsWith(rod fyne.CanvasObject) bool {
	return d.internalFind(d.Position(), rod)
}

func (d *Disk) GetCurrentRod(rods []fyne.CanvasObject) int { // TODO: так себе метод
	for i, r := range rods {
		if d.internalFind(d.boundPosition, r) {
			return i
		}
	}
	return -1
}

func (d *Disk) internalFind(position fyne.Position, rod fyne.CanvasObject) bool {
	r1Left := rod.Position().X
	r1Right := rod.Position().X + rod.Size().Width
	r1Top := rod.Position().Y
	r1Bottom := rod.Position().Y + rod.Size().Height

	r2Left := position.X
	r2Right := position.X + d.Size().Width
	r2Top := position.Y
	r2Bottom := position.Y + d.Size().Height

	return r1Right > r2Left && r1Left < r2Right && r1Bottom > r2Top && r1Top < r2Bottom
}

func (d *Disk) Release() {
	d.blocked = true
}
