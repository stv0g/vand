// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package widget

import (
	"encoding/json"

	"github.com/stv0g/vand/pkg/store"
	"github.com/tdewolff/canvas"
)

type Gauge struct {
	WidgetBase

	Min float64 `yaml:"min"`
	Max float64 `yaml:"min"`

	Value   float64 `yaml:"value"`
	ValueOf string  `yaml:"value_of"`

	Position [2]float64 `yaml:"position"`
	Size     [2]float64 `yaml:"size"`
}

func (w *Gauge) UnmarshalJSON(b []byte) error {
	type gauge Gauge
	if err := json.Unmarshal(b, (*Gauge)(w)); err != nil {
		return err
	}

	return nil
}

func (w *Gauge) Init() error {
	return nil
}

func (W *Gauge) Close() error {
	return nil
}

func (w *Gauge) Draw(ctx *canvas.Context, s *store.Store) error {
	// ctx.DrawImage(w.rect.X, w.rect.Y, w.image, w.resolution)

	return nil
}
