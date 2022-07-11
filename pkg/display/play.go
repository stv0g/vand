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
	var current_page *Page

	dim := Pixels / DotsPerMillimeter

	for _, page := range pages {
		// Get first page
		if current_page == nil {
			current_page = page
		}

		if err := page.Init(); err != nil {
			return fmt.Errorf("failed to initialize page: %w", err)
		}
	}

	c := canvas.New(dim, dim)

	for current_page != nil {
		c.Reset()
		ctx := canvas.NewContext(c)

		// Draw background
		if current_page.BackgroundColor != nil {
			ctx.SetFillColor(current_page.BackgroundColor)
			ctx.DrawPath(0, 0, canvas.Rectangle(Pixels, Pixels))
		}

		if err := current_page.Draw(ctx, s); err != nil {
			return fmt.Errorf("failed to draw page: %w", err)
		}

		rst := rasterizer.Draw(c, Resolution, canvas.LinearColorSpace{})
		d.Draw(rst.Bounds(), rst, image.ZP)

		t := time.NewTimer(current_page.Time)
		<-t.C

		current_page = current_page.Next
	}

	return nil
}
