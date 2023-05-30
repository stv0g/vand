// Copyright 2016 The Periph Authors. All rights reserved.
// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package ssd1351

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/stv0g/vand/pkg/devices/display/ssd1351/image16bit"
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/display"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
)

// NewSPI returns a Dev object that communicates over SPI to a SSD1351 display
// controller.
//
// The SSD1351 can operate at up to 3.3Mhz, which is much higher than I²C. This
// permits higher refresh rates.
//
// # Wiring
//
// Connect SDA to SPI_MOSI, SCK to SPI_CLK, CS to SPI_CS.
func NewSPI(p spi.Port, dc gpio.PinOut, rst gpio.PinOut) (*Dev, error) {
	c, err := p.Connect(32*physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}

	return newDev(c, dc, rst)
}

// Dev is an open handle to the display controller.
type Dev struct {
	// Communication
	c   conn.Conn
	dc  gpio.PinOut
	rst gpio.PinOut

	// Display size controlled by the SSD1351.
	rect image.Rectangle

	// Mutable
	// See page 25 for the GDDRAM pages structure.
	// Narrow screen will waste the end of each page.
	// Short screen will ignore the lower pages.
	// There is 8 pages, each covering an horizontal band of 8 pixels high (1
	// byte) for 128 bytes.
	// 8*128 = 1024 bytes total for 128x64 display.
	buffer []byte
	// next is lazy initialized on first Draw(). Write() skips this buffer.
	next   *image16bit.VerticalLSB
	halted bool
}

func (d *Dev) String() string {
	return fmt.Sprintf("ssd1360.Dev{%s, %s, %s}", d.c, d.dc, d.rect.Max)
}

// ColorModel implements display.Drawer.
//
// It is a one bit color model, as implemented by image16bit.Bit.
func (d *Dev) ColorModel() color.Model {
	return image16bit.BitsModel
}

// Bounds implements display.Drawer. Min is guaranteed to be {0, 0}.
func (d *Dev) Bounds() image.Rectangle {
	return d.rect
}

// Draw implements display.Drawer.
//
// It draws synchronously, once this function returns, the display is updated.
// It means that on slow bus (I²C), it may be preferable to defer Draw() calls
// to a background goroutine.
func (d *Dev) Draw(r image.Rectangle, src image.Image, sp image.Point) error {
	var next []byte
	if img, ok := src.(*image16bit.VerticalLSB); ok && r == d.rect && img.Rect == d.rect && sp.X == 0 && sp.Y == 0 {
		// Exact size, full frame, image16bit encoding: fast path!
		next = img.Pix
	} else {
		// Double buffering.
		if d.next == nil {
			d.next = image16bit.NewVerticalLSB(d.rect)
		}
		next = d.next.Pix
		draw.Src.Draw(d.next, r, src, sp)
	}
	return d.drawInternal(next)
}

// Write writes a buffer of pixels to the display.
//
// This function accepts the content of image16bit.VerticalLSB.Pix.
func (d *Dev) Write(pixels []byte) (int, error) {
	if len(pixels) != len(d.buffer) {
		return 0, fmt.Errorf("ssd1351: invalid pixel stream length; expected %d bytes, got %d bytes", len(d.buffer), len(pixels))
	}

	// Write() skips d.next so it saves 1kb of RAM.
	if err := d.drawInternal(pixels); err != nil {
		return 0, err
	}

	return len(pixels), nil
}

// Invert the display (black on white vs white on black).
func (d *Dev) Invert(blackOnWhite bool) error {
	b := []byte{0xa6}
	if blackOnWhite {
		b[0] = 0xa7
	}
	return d.sendCommand(b)
}

// Halt turns off the display.
//
// Sending any other command afterward reenables the display.
func (d *Dev) Halt() error {
	d.halted = false
	err := d.sendCommand([]byte{0xae})
	if err == nil {
		d.halted = true
	}
	return err
}

// newDev is the common initialization code that is independent of the
// communication protocol (I²C or SPI) being used.
func newDev(c conn.Conn, dc gpio.PinOut, rst gpio.PinOut) (*Dev, error) {
	d := &Dev{
		c:      c,
		dc:     dc,
		rst:    rst,
		rect:   image.Rect(0, 0, 128, 128),
		buffer: make([]byte, 128*128*2),
	}

	if err := d.reset(); err != nil {
		return nil, err
	}

	if err := d.init(); err != nil {
		return nil, err
	}

	time.Sleep(200 * time.Millisecond)

	if err := d.sendCommand([]byte{0xaf}); err != nil { // Turn on the OLED display
		return nil, fmt.Errorf("failed to turn on display: %w", err)
	}

	return d, nil
}

func (d *Dev) reset() error {
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)
	if err := d.rst.Out(gpio.Low); err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)
	if err := d.rst.Out(gpio.High); err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)

	return nil
}

func (d *Dev) init() error {
	if err := d.sendCommand([]byte{0xfd}); err != nil { // Set Command Lock
		return err
	}
	if err := d.sendData([]byte{0x12}); err != nil { // Unlock OLED driver IC MCU interface from entering command [reset]
		return err
	}

	if err := d.sendCommand([]byte{0xfd}); err != nil { // Set Command Lock
		return err
	}
	if err := d.sendData([]byte{0xb1}); err != nil { // Command A2,B1,B3,BB,BE,C1 accessible if in unlock state
		return err
	}

	if err := d.sendCommand([]byte{0xae}); err != nil { // Set Sleep Mode: ON
		return err
	}
	if err := d.sendData([]byte{0xa4}); err != nil { // Set Display Mode: All Off
		return err
	}

	if err := d.sendCommand([]byte{0x15}); err != nil { // Set Column Address
		return err
	}
	if err := d.sendData([]byte{0x00, 0x7f}); err != nil { // Start, End
		return err
	}

	if err := d.sendCommand([]byte{0x75}); err != nil { // Set Row Address
		return err
	}
	if err := d.sendData([]byte{0x00, 0x7f}); err != nil { // Start, End
		return err
	}

	if err := d.sendCommand([]byte{0xb3}); err != nil { // Set Front Clock Divider / Oscillator Frequency
		return err
	}
	if err := d.sendData([]byte{0xf1}); err != nil {
		return err
	}

	if err := d.sendCommand([]byte{0xca}); err != nil { // Set Multiplex Ratio
		return err
	}
	if err := d.sendData([]byte{0x7f}); err != nil {
		return err
	}

	if err := d.sendCommand([]byte{0xa0}); err != nil { // Set Re-map & Dual COM Line Mode
		return err
	}
	if err := d.sendData([]byte{0x74}); err != nil { // Horizontal address increment
		return err
	}

	if err := d.sendCommand([]byte{0xa1}); err != nil { // Set Display Start Line
		return err
	}
	if err := d.sendData([]byte{0x00}); err != nil { // start 00 line
		return err
	}

	if err := d.sendCommand([]byte{0xa2}); err != nil { // Set Display Offset
		return err
	}
	if err := d.sendData([]byte{0x00}); err != nil {
		return err
	}

	if err := d.sendCommand([]byte{0xab, 0x01}); err != nil { // Set Function selection
		return err
	}

	if err := d.sendCommand([]byte{0xb4}); err != nil { // Set Segment Low Voltage (VSL)
		return err
	}
	if err := d.sendData([]byte{0xa0, 0xb5, 0x55}); err != nil {
		return err
	}

	if err := d.sendCommand([]byte{0xc1}); err != nil { // Set Contrast Current for Color A,B,C
		return err
	}
	if err := d.sendData([]byte{0xc8, 0x80, 0xc0}); err != nil {
		return err
	}

	if err := d.sendCommand([]byte{0xc7}); err != nil { // Master Contrast Current Control
		return err
	}
	if err := d.sendData([]byte{0x0f}); err != nil {
		return err
	}

	if err := d.sendCommand([]byte{0xb1}); err != nil { // Set Phase Length
		return err
	}
	if err := d.sendData([]byte{0x32}); err != nil {
		return err
	}

	if err := d.sendCommand([]byte{0xb2}); err != nil { // Display Enhancement
		return err
	}
	if err := d.sendData([]byte{0xa4, 0x00, 0x00}); err != nil {
		return err
	}

	if err := d.sendCommand([]byte{0xbb}); err != nil { // Set Pre-charge voltage
		return err
	}
	if err := d.sendData([]byte{0x17}); err != nil {
		return err
	}

	if err := d.sendCommand([]byte{0xb6}); err != nil { // Set Second Pre-charge period
		return err
	}
	if err := d.sendData([]byte{0x01}); err != nil {
		return err
	}

	if err := d.sendCommand([]byte{0xbe}); err != nil { // Set V_COMH Voltage
		return err
	}
	if err := d.sendData([]byte{0x05}); err != nil {
		return err
	}

	if err := d.sendCommand([]byte{0xa6}); err != nil { // Set Display Mode: Reset to normal display [reset]
		return err
	}

	return nil
}

// drawInternal sends image data to the controller.
func (d *Dev) drawInternal(next []byte) error {
	if err := d.sendCommand([]byte{0x5c}); err != nil {
		return err
	}

	if err := d.sendData(next); err != nil {
		return err
	}

	return nil
}

func (d *Dev) sendData(c []byte) error {
	if d.halted {
		// Transparently enable the display.
		if err := d.sendCommand(nil); err != nil {
			return err
		}
	}
	// 4-wire SPI.
	if err := d.dc.Out(gpio.High); err != nil {
		return err
	}

	maxTxSz := d.c.(conn.Limits).MaxTxSize()

	sz := 0
	for e := c; len(e) > 0; e = e[sz:] {
		sz = len(e)
		if sz > maxTxSz {
			sz = maxTxSz
		}

		if err := d.c.Tx(e[:sz], nil); err != nil {
			return err
		}
	}

	return nil
}

func (d *Dev) sendCommand(c []byte) error {
	if d.halted {
		// Transparently enable the display.
		c = append([]byte{0xaf}, c...)
		d.halted = false
	}
	// 4-wire SPI.
	if err := d.dc.Out(gpio.Low); err != nil {
		return err
	}
	return d.c.Tx(c, nil)
}

var _ display.Drawer = &Dev{}
