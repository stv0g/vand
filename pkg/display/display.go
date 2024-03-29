// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package display

import (
	"github.com/tdewolff/canvas"
	"periph.io/x/conn/v3/display"
)

const (
	DotsPerMillimeter = 1.0 / 0.21

	Pixels = 128
)

var Resolution = canvas.DPMM(DotsPerMillimeter)

type Display struct {
	display.Drawer

	next chan struct{}
}
