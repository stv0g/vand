package display

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/stv0g/vand/pkg/store"
	"github.com/tdewolff/canvas"
)

type Widget interface {
	io.Closer

	Init() error
	Draw(ctx *canvas.Context, s *store.Store) error
}

type WidgetBase struct {
	Type string `json:"type"`
}

func unmarshalWidgets(raws []json.RawMessage) ([]Widget, error) {
	widgets := []Widget{}

	for _, raw := range raws {
		var w WidgetBase

		if err := json.Unmarshal(raw, &w); err != nil {
			return nil, err
		}

		var i Widget
		switch w.Type {
		case "text":
			i = &TextWidget{}
		case "image":
			i = &ImageWidget{}
		default:
			return nil, errors.New("unknown widget type")
		}

		if err := json.Unmarshal(raw, i); err != nil {
			return nil, err
		}

		widgets = append(widgets, i)
	}

	return widgets, nil
}
