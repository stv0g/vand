package display

import (
	"fmt"
	"image"
	"time"

	"github.com/stv0g/vand/pkg/store"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/rasterizer"
)

func (d *Display) Play(pages map[string]*Page, s *store.Store) error {
	dim := Pixels / DotsPerMillimeter
	// dim = 500

	for _, page := range pages {
		if err := page.Init(); err != nil {
			return fmt.Errorf("failed to initialize page: %w", err)
		}
	}

	c := canvas.New(dim, dim)

	for {
		for _, page := range pages {
			c.Reset()
			ctx := canvas.NewContext(c)

			if err := page.Draw(ctx, s); err != nil {
				return fmt.Errorf("failed to draw page: %w", err)
			}

			rst := rasterizer.Draw(c, Resolution, canvas.SRGBColorSpace{})
			d.Draw(rst.Bounds(), rst, image.ZP)

			t := time.NewTimer(page.Time)
			<-t.C
		}
	}

	return nil
}
