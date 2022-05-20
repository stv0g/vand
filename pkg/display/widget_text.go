package display

import (
	"bytes"
	"fmt"
	"image/color"
	"text/template"

	"github.com/stv0g/vand/pkg/store"
	"github.com/tdewolff/canvas"
)

type TextWidget struct {
	WidgetBase

	Size      float32
	Color     color.Color
	FontFace  *canvas.FontFace
	TextAlign canvas.TextAlign

	Template template.Template
}

func (w *TextWidget) Init() error {
	return nil
}

func (W *TextWidget) Close() error {
	return nil
}

func (w *TextWidget) Draw(ctx *canvas.Context, s *store.Store) error {
	wr := &bytes.Buffer{}

	data := s.GetValues()

	if err := w.Template.Execute(wr, data); err != nil {
		return fmt.Errorf("failed to get values from store: %w", err)
	}

	ctx.SetStrokeColor(w.Color)

	// txt := canvas.NewTextLine(w.FontFace, wr.String(), w.TextAlign)
	// ctx.DrawText(w.Position.X, w.Position.Y, txt)

	return nil
}
