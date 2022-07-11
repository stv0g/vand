package display

import (
	"fmt"

	"github.com/stv0g/vand/pkg/config"
	"github.com/stv0g/vand/pkg/store"
	"github.com/tdewolff/canvas"
)

type Page struct {
	config.DisplayPage

	Over *Page
	Next *Page
}

func (p *Page) Draw(c *canvas.Context, s *store.Store) error {
	if p.Over != nil {
		if err := p.Over.Draw(c, s); err != nil {
			return err
		}
	}

	for _, widget := range p.Widgets {
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

func LoadPages(configPages []config.DisplayPage) (map[string]*Page, error) {
	pages := map[string]*Page{}
	page_ids := []string{}
	for _, page := range configPages {
		pages[page.ID] = &Page{DisplayPage: page}
		page_ids = append(page_ids, page.ID)
	}

	i := 0
	for _, page := range pages {
		if page.NextID == "" {
			page.NextID = page_ids[(i+1)%len(pages)]
		}

		var ok bool
		page.Next, ok = pages[page.NextID]
		if !ok {
			return nil, fmt.Errorf("could not find page with id: %s", page.NextID)
		}

		if page.OverID != "" {
			page.Over, ok = pages[page.OverID]
			if !ok {
				return nil, fmt.Errorf("could not find page with id: %s", page.OverID)
			}
		}

		i++
	}

	return pages, nil
}
