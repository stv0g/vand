// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

//go:build virtual

package display

import (
	"log"

	"github.com/stv0g/vand/pkg/devices/display/virtual"
)

func NewVirtualDisplay() (*Display, error) {
	dev, err := virtual.New()
	if err != nil {
		log.Fatalf("Failed to open display: %s", err)
	}

	return &Display{
		Drawer: dev,
	}, nil
}
