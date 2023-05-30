// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package display

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"time"

	"github.com/nfnt/resize"
)

// convertAndResizeAndCenter takes an image, resizes and centers it on a
// image.Gray of size w*h.
func convertAndResizeAndCenter(w, h int, src image.Image) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	src = resize.Thumbnail(uint(w), uint(h), src, resize.Bicubic)
	r := src.Bounds()
	r = r.Add(image.Point{(w - r.Max.X) / 2, (h - r.Max.Y) / 2})
	draw.Draw(img, r, src, image.Point{}, draw.Src)

	return img
}

func (d *Display) PlayGIF(rd io.Reader) error {
	g, err := gif.DecodeAll(rd)
	if err != nil {
		return fmt.Errorf("failed to decode gif: %w", err)
	}

	// Converts every frame to image.Gray and resize them:
	imgs := make([]*image.RGBA, len(g.Image))

	overpaintImg := image.NewRGBA(image.Rect(0, 0, g.Config.Width, g.Config.Height))

	draw.Draw(overpaintImg, overpaintImg.Bounds(), g.Image[0], image.Point{}, draw.Src)
	for i := range g.Image {
		draw.Draw(overpaintImg, overpaintImg.Bounds(), g.Image[i], image.Point{}, draw.Over)

		imgs[i] = convertAndResizeAndCenter(128, 128, overpaintImg)
	}

	// Display the frames in a loop:
	for i := 0; ; i++ {
		index := i % len(imgs)
		c := time.After(time.Duration(10*g.Delay[index]) * time.Millisecond)
		img := imgs[index]
		if err := d.Draw(img.Bounds(), img, image.Point{}); err != nil {
			return err
		}
		<-c
	}
}
