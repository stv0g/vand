package display

import (
	"log"

	"github.com/stv0g/vand/pkg/config"
	"github.com/stv0g/vand/pkg/devices/display/mockup"
)

func NewMockupDisplay(cfg *config.Display) (*Display, error) {
	dev, err := mockup.New()
	if err != nil {
		log.Fatalf("Failed to open display: %s", err)
	}

	return &Display{
		Drawer: dev,
	}, nil
}
