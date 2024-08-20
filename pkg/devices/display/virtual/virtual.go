// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

//nolint:typecheck
package virtual

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/gio"
)

const (
	Width  = 128
	Height = 128
)

type Dev struct {
	window *app.Window
	rect   image.Rectangle
	image  draw.Image
}

func New() (*Dev, error) {
	w := &app.Window{}
	w.Option(
		app.Title("vand"),
		app.Size(unit.Dp(Width), unit.Dp(Height)),
		app.MinSize(unit.Dp(Width), unit.Dp(Height)),
		app.MaxSize(unit.Dp(Width), unit.Dp(Height)))

	d := &Dev{
		window: w,
		rect: image.Rectangle{
			Max: image.Point{
				Width,
				Height,
			},
		},
	}

	d.image = image.NewRGBA(d.rect)

	go d.loop()

	return d, nil
}

func (d *Dev) loop() {
	var ops op.Ops
	for {
		e := d.window.Event()

		switch e := e.(type) {

		case app.DestroyEvent:
			os.Exit(0)

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				c := gio.NewContain(gtx, float64(d.rect.Dx()), float64(d.rect.Dy()))
				c.RenderImage(d.image, canvas.Identity)
				return c.Dimensions()
			})

			e.Frame(gtx.Ops)
		}
	}
}

func (d *Dev) String() string {
	return fmt.Sprintf("virtual.Dev{%s}", d.rect.Max)
}

// ColorModel implements display.Drawer.
func (d *Dev) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds implements display.Drawer. Min is guaranteed to be {0, 0}.
func (d *Dev) Bounds() image.Rectangle {
	return d.rect
}

// Draw implements display.Drawer.
//
// It draws synchronously, once this function returns, the display is updated.
// It means that on slow bus (I²C), it may be preferable to defer Draw() calls
// to a background goroutine.
func (d *Dev) Draw(r image.Rectangle, src image.Image, sp image.Point) error {
	draw.Draw(d.image, r, src, sp, draw.Src)
	d.window.Invalidate()
	return nil
}

func (d *Dev) Halt() error {
	log.Println("Halt not supported yet")

	return nil
}
