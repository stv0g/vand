// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package widget

import (
	"encoding/json"

	"github.com/stv0g/vand/pkg/store"
	"github.com/tdewolff/canvas"
)

type Line struct {
	WidgetBase

	From [2]float64 `yaml:"from"`
	To   [2]float64 `yaml:"to"`
}

func (w *Line) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, (*Line)(w))
}

func (w *Line) Init() error {
	return nil
}

func (W *Line) Close() error {
	return nil
}

func (w *Line) Draw(ctx *canvas.Context, s *store.Store) error {
	ctx.MoveTo(w.From[0], w.From[1])
	ctx.LineTo(w.To[0], w.To[1])

	return nil
}
