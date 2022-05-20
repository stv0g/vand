package display

import (
	"encoding/json"
	"image/color"
	"time"

	"github.com/stv0g/vand/pkg/store"
	"github.com/tdewolff/canvas"
)

type Page struct {
	ID   string `json:"id"`
	Next string `json:"next"`
	Over string `json:"over"`

	Time            time.Duration `json:"time"`
	BackgroundColor color.Color   `json:"background-color"`

	Widgets    []Widget          `json:"-"`
	WidgetsRaw []json.RawMessage `json:"widgets"`
}

func (p *Page) UnmarshalJSON(b []byte) error {
	type page Page

	if err := json.Unmarshal(b, (*page)(p)); err != nil {
		return err
	}

	var err error
	if p.Widgets, err = unmarshalWidgets(p.WidgetsRaw); err != nil {
		return err
	}

	return nil
}

func (p *Page) Draw(c *canvas.Context, s *store.Store) error {

	for _, widget := range p.Widgets {
		if err := widget.Draw(c, s); err != nil {
			return err
		}
	}

	return nil
}
