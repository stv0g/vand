// Copyright 2016 The Periph Authors. All rights reserved.
// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

// Package image16bit implements 16 bit (5, 6, 7 bits per color) 2D graphics.
//
// It is compatible with package image/draw.
package image16bit

import (
	"encoding/binary"
	"image"
	"image/color"
	"image/draw"
)

// Bit implements a 65k color.
type Bits uint16

const (
	Black = 0
	White = 0xffff
)

func NewBits(c color.Color) Bits {
	r, g, b, _ := c.RGBA()

	return Bits((r>>11)<<11 | (g>>10)<<5 | (b >> 11))
}

func (bi Bits) RGBA() (uint32, uint32, uint32, uint32) {
	i := uint32(bi)

	r := ((i >> 11) & 0x1f) << 11
	g := ((i >> 5) & 0x3f) << 10
	b := ((i >> 0) & 0x1f) << 11

	return r, g, b, 0xffff
}

var BitsModel = color.ModelFunc(convert)

// VerticalLSB is a 16 bit RGB image.
//
// # Each 2 bytes represent a single pixel
//
// It is designed specifically to work with SSD1351 OLED display controller.
type VerticalLSB struct {
	// Pix holds the image's pixels, as vertically LSB-first packed bitmap. It
	// can be passed directly to ssd1351.Dev.Write()
	Pix []byte

	Stride int

	// Rect is the image's bounds.
	Rect image.Rectangle
}

// NewVerticalLSB returns an initialized VerticalLSB instance.
func NewVerticalLSB(r image.Rectangle) *VerticalLSB {
	sz := r.Size()

	return &VerticalLSB{
		Pix:    make([]byte, sz.X*sz.Y*2),
		Rect:   r,
		Stride: r.Dx(),
	}
}

// ColorModel implements image.Image.
func (i *VerticalLSB) ColorModel() color.Model {
	return BitsModel
}

// Bounds implements image.Image.
func (i *VerticalLSB) Bounds() image.Rectangle {
	return i.Rect
}

// At implements image.Image.
func (i *VerticalLSB) At(x, y int) color.Color {
	return i.BitsAt(x, y)
}

// BitAt is the optimized version of At().
func (i *VerticalLSB) BitsAt(x, y int) Bits {
	if !(image.Point{x, y}.In(i.Rect)) {
		return 0
	}

	offset := i.PixOffset(x, y)
	return Bits(binary.BigEndian.Uint16(i.Pix[offset:]))
}

func (i *VerticalLSB) PixOffset(x, y int) int {
	pY := y - i.Rect.Min.Y
	offset := pY*i.Stride + (x - i.Rect.Min.X)
	return offset * 2
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (i *VerticalLSB) Opaque() bool {
	return true
}

// Set implements draw.Image
func (i *VerticalLSB) Set(x, y int, c color.Color) {
	i.SetBit(x, y, convertBit(c))
}

// SetBit is the optimized version of Set().
func (i *VerticalLSB) SetBit(x, y int, b Bits) {
	if !(image.Point{x, y}.In(i.Rect)) {
		return
	}
	offset := i.PixOffset(x, y)
	binary.BigEndian.PutUint16(i.Pix[offset:], uint16(b))
}

var _ draw.Image = &VerticalLSB{}

// Anything not transparent and not pure black is white.
func convert(c color.Color) color.Color {
	return convertBit(c)
}

// Anything not transparent and not pure black is white.
func convertBit(c color.Color) Bits {
	switch t := c.(type) {
	case Bits:
		return t
	default:
		return NewBits(c)
	}
}
