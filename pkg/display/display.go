package display

import (
	"github.com/tdewolff/canvas"
	"periph.io/x/conn/v3/display"
)

const (
	DotsPerMillimeter = 1.0 / 0.21

	Pixels = 128
)

var (
	Resolution = canvas.DPMM(DotsPerMillimeter)
)

type Display struct {
	display.Drawer

	Pages map[string]Page

	next chan struct{}
}
