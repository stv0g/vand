package display

import (
	"image"
	"time"

	"github.com/stv0g/vand/pkg/store"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/rasterizer"
)

func (d *Display) PlaybackPages(pages []Page, s *store.Store) error {
	dim := Pixels / DotsPerMillimeter
	// dim = 500

	c := canvas.New(dim, dim)

	for _, page := range pages {
		c.Reset()
		ctx := canvas.NewContext(c)

		if err := page.Draw(ctx, s); err != nil {

		}

		rst := rasterizer.Draw(c, Resolution, canvas.SRGBColorSpace{})
		d.Draw(rst.Bounds(), rst, image.ZP)

		t := time.NewTimer(0)
		<-t.C
	}

	return nil
}
