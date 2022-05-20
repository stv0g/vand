package display

import (
	"io/ioutil"

	"github.com/tdewolff/canvas"

	"github.com/stv0g/vand/resources"
)

var fontFamily *canvas.FontFamily

func init() {
	fontFamily = canvas.NewFontFamily("FireCode")

	f, err := resources.FS.Open("fonts/FiraCode-Regular.ttf")
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(f)

	if err := fontFamily.LoadFont(b, 0, canvas.FontRegular); err != nil {
		panic(err)
	}
}

func Draw(ctx *canvas.Context) error {
	ctx.SetFillColor(canvas.Black)
	ctx.DrawPath(0, 0, canvas.Rectangle(ctx.Width(), ctx.Height()))

	ctx.SetFillColor(canvas.Black)

	// lenna, err := resources.FS.Open("lenna.png")
	// if err != nil {
	// 	return err
	// }

	// img, err := png.Decode(lenna)
	// if err != nil {
	// 	return err
	// }

	// ctx.DrawImage(0, 0, img, canvas.DPMM(15.0))

	headerFace := fontFamily.Face(8.0, canvas.Green, canvas.FontRegular, canvas.FontNormal)
	ctx.DrawText(0, ctx.Height()/2, canvas.NewTextBox(headerFace, "Hello World", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))

	ctx.SetFillColor(canvas.Red)
	ctx.SetStrokeWidth(1.0)
	ctx.SetStrokeColor(canvas.Red)
	ctx.DrawPath(ctx.Width()/2, ctx.Height()/2, canvas.Circle(10))

	return nil
}
