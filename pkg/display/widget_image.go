package display

import (
	"encoding/json"
	"image"

	"github.com/stv0g/vand/pkg/store"
	"github.com/tdewolff/canvas"
)

type ImageWidget struct {
	WidgetBase

	File string `json:"file"`
	// Animate bool   `json:"animate"`
	// Loop    bool   `json:"loop"`

	Position   [2]float64        `json:"position"`
	Size       [2]float64        `json:"size"`
	Resolution canvas.Resolution `json:"resolution"`

	rect  canvas.Rect `json:"-"`
	image image.Image `json:"-"`
}

func (w *ImageWidget) UnmarshalJSON(b []byte) error {
	type imageWidget ImageWidget
	if err := json.Unmarshal(b, (*imageWidget)(w)); err != nil {
		return err
	}

	w.rect = canvas.Rect{
		X: w.Position[0],
		Y: w.Position[1],
		W: w.Size[0],
		H: w.Size[1],
	}

	return nil
}

func (w *ImageWidget) Init() error {
	return nil
}

func (W *ImageWidget) Close() error {
	return nil
}

func (w *ImageWidget) Draw(ctx *canvas.Context, s *store.Store) error {
	// ctx.DrawImage(w.rect.X, w.rect.Y, w.image, w.resolution)

	return nil
}
