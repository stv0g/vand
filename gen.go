// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

//go:build ignore
// +build ignore

// This program generates font7x13.go.
//
// It exists so cmd/ssd1306 does not depend on golang.org/x/image/...
//
// This program is not built by default.

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"image"
	"io/ioutil"
	"os"
	"text/template"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"periph.io/x/devices/v3/ssd1306/image1bit"
)

var text = `// Code generated by "go run gen.go"; DO NOT EDIT.
package main
// This data is derived from files in the font/fixed directory of the Plan 9
// Port source code (https://github.com/9fans/plan9port) which were originally
// based on the public domain X11 misc-fixed font files.
import (
	"image"
	"image/color"
	"image/draw"
	"periph.io/x/devices/v3/ssd1306/image1bit"
)
type bit bool
func (b bit) RGBA() (uint32, uint32, uint32, uint32) {
  if b {
    return 65535, 65535, 65535, 65535
  }
  return 0, 0, 0, 0
}
func convertBit(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return bit((r | g | b) >= 0x8000)
}
type alpha struct {
	image1bit.VerticalLSB
}
func (a *alpha) ColorModel() color.Model {
	return color.ModelFunc(convertBit)
}
func (a *alpha) At(x, y int) color.Color {
	return convertBit(a.VerticalLSB.At(x, y))
}
// glyphs contains chars 0x21 to 0x7F.
var glyphs = []alpha{
{{range .}}	{
	image1bit.VerticalLSB{
			Pix:    []byte{ {{range .Pix}}{{.}}, {{end}} },
			Stride: {{.Stride}},
			Rect:   image.Rectangle{Max: image.Point{ {{.Rect.Max.X}}, {{.Rect.Max.Y}} } },
		},
	},
{{end}}
}
// drawText draws text on an image.
//
// It is intentionally very limited. Use golang.org/x/image/font for complete
// functionality.
func drawText(dst draw.Image, p image.Point, t string) {
	const base = 0x21
	r := image.Rect(0, 0, 7, 13).Add(p)
	u := image.Uniform{C: image1bit.On}
	for _, c := range t {
		if c >= base && int(c-base) < len(glyphs) {
			draw.DrawMask(dst, r, &u, image.Point{}, &glyphs[c-base], image.Point{}, draw.Over)
		}
		r = r.Add(image.Point{7, 0})
	}
}
`

func mainImpl() error {
	t, err := template.New("main").Parse(text)
	if err != nil {
		return err
	}
	const base = 0x21
	glyphs := [0x80 - base]image1bit.VerticalLSB{}
	for i := range glyphs {
		glyphs[i] = *image1bit.NewVerticalLSB(image.Rect(0, 0, 6, 13))
		drawer := font.Drawer{
			Src:  &image.Uniform{C: image1bit.On},
			Dst:  &glyphs[i],
			Face: basicfont.Face7x13,
			Dot:  fixed.P(0, 12),
		}
		drawer.DrawString(string(rune(i + base)))
	}

	var b bytes.Buffer
	if err = t.Execute(&b, glyphs); err != nil {
		return err
	}
	src, err := format.Source(b.Bytes())
	if err != nil {
		return err
	}
	return ioutil.WriteFile("font7x13.go", src, 0644)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "gen: %s.\n", err)
		os.Exit(1)
	}
}