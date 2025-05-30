// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package widget

import (
	"encoding/json"

	"github.com/stv0g/vand/pkg/store"
	"github.com/tdewolff/canvas"
)

type Image struct {
	WidgetBase

	File       string            `yaml:"file"`
	Animate    bool              `yaml:"animate"`
	Loop       bool              `yaml:"loop"`
	Position   [2]float64        `yaml:"position"`
	Size       [2]float64        `yaml:"size"`
	Resolution canvas.Resolution `yaml:"resolution"`

	rect canvas.Rect
}

func (w *Image) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, (*Image)(w)); err != nil {
		return err
	}

	w.rect = canvas.Rect{
		X0: w.Position[0],
		Y0: w.Position[1],
		X1: w.Position[0] + w.Size[0],
		Y1: w.Position[1] + w.Size[1],
	}

	return nil
}

func (w *Image) Init() error {
	return nil
}

func (W *Image) Close() error {
	return nil
}

func (w *Image) Draw(ctx *canvas.Context, s *store.Store) error {
	// ctx.DrawImage(w.rect.X, w.rect.Y, w.image, w.resolution)

	return nil
}
