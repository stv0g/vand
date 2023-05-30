// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package widget

import (
	"bytes"
	"fmt"
	"image/color"
	"text/template"

	"github.com/stv0g/vand/pkg/store"
	"github.com/tdewolff/canvas"
)

type Text struct {
	WidgetBase

	Position [2]float64  `yaml:"position"`
	Size     float64     `yaml:"size"`
	Color    color.Color `yaml:"color"`
	Font     string      `yaml:"font"`
	Template string      `yaml:"template"`
	Align    string      `yaml:"align"`

	fontFace  *canvas.FontFace
	textAlign canvas.TextAlign
	template  *template.Template
}

func (w *Text) Init() error {
	w.template = template.New("widget")
	if _, err := w.template.Parse(w.Template); err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	ff := canvas.NewFontFamily("ff")
	ff.LoadLocalFont(w.Font, canvas.FontRegular)

	w.fontFace = ff.Face(w.Size, w.Color, canvas.FontRegular, canvas.FontNormal)

	w.textAlign = canvas.Left

	return nil
}

func (W *Text) Close() error {
	return nil
}

func (w *Text) Draw(ctx *canvas.Context, s *store.Store) error {
	wr := &bytes.Buffer{}

	data := s.Flatten(".")

	if err := w.template.Execute(wr, data); err != nil {
		return fmt.Errorf("failed to get values from store: %w", err)
	}

	ctx.SetStrokeColor(w.Color)

	txt := canvas.NewTextLine(w.fontFace, wr.String(), w.textAlign)
	ctx.DrawText(w.Position[0], w.Position[1], txt)

	return nil
}
