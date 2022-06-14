package display

import (
	"fmt"
	"log"

	"github.com/stv0g/vand/pkg/config"
	"github.com/stv0g/vand/pkg/store"
	"github.com/tdewolff/canvas"
)

type Page struct {
	config.DisplayPage

	over *Page
	next *Page
}

func (p *Page) Draw(c *canvas.Context, s *store.Store) error {
	if p.over != nil {
		if err := p.over.Draw(c, s); err != nil {
			return err
		}
	}

	for _, widget := range p.Widgets {
		log.Printf("draw widget %+#v", widget)
		if err := widget.Draw(c, s); err != nil {
			return err
		}
	}

	return nil
}

func (p *Page) Init() error {
	for _, widget := range p.Widgets {
		if err := widget.Init(); err != nil {
			return fmt.Errorf("failed to initialize widget: %w", err)
		}
	}

	return nil
}
