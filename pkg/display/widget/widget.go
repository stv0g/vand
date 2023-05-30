// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package widget

import (
	"io"

	"github.com/stv0g/vand/pkg/store"
	"github.com/tdewolff/canvas"
)

type Widget interface {
	io.Closer

	Init() error
	Draw(ctx *canvas.Context, s *store.Store) error
}

type WidgetBase struct {
	Type string `yaml:"type"`
}
