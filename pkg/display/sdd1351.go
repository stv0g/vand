// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package display

import (
	"fmt"
	"log"
	"time"

	"github.com/stv0g/vand/pkg/config"
	"github.com/stv0g/vand/pkg/devices/display/ssd1351"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

func NewDisplay(cfg *config.Display) (*Display, error) {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	spiPort, err := spireg.Open(cfg.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to open SPI port")
	}
	// defer spiPort.Close()

	dcPin := gpioreg.ByName(cfg.Pins.DC)
	if dcPin == nil {
		return nil, fmt.Errorf("failed to find dc pin")
	}

	resetPin := gpioreg.ByName(cfg.Pins.Reset)
	if resetPin == nil {
		return nil, fmt.Errorf("failed to find reset pin")
	}

	nextPin := gpioreg.ByName(cfg.Pins.Next)
	if nextPin != nil {
		if err := nextPin.In(gpio.PullUp, gpio.FallingEdge); err != nil {
			return nil, fmt.Errorf("failed to setup input pin")
		}
	}

	dev, err := ssd1351.NewSPI(spiPort, dcPin, resetPin)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize display")
	}

	return &Display{
		Drawer: dev,
		next:   waitForEdge(nextPin),
	}, nil
}

func waitForEdge(pin gpio.PinIO) chan struct{} {
	if pin == nil {
		return nil
	}

	ch := make(chan struct{})

	go func() {
		for {
			if pin.WaitForEdge(-1) {
				ch <- struct{}{}
			}
		}
	}()

	return ch
}
