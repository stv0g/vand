package widget

import (
	"fmt"
	"image/color"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/image/colornames"
)

// Decode takes an input structure and uses reflection to translate it to
// the output structure. output must be a pointer to a map or struct.
func decodeWithHooks(input interface{}, output interface{}, hook mapstructure.DecodeHookFunc) error {
	config := &mapstructure.DecoderConfig{
		Metadata:   nil,
		Result:     output,
		DecodeHook: hook,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func DecodeHookFunc() mapstructure.DecodeHookFunc {
	return mapstructure.ComposeDecodeHookFunc(
		widgetDecodeHook,
		colorDecodeHook,
	)
}

func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff

	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}

	return
}

func colorDecodeHook(
	f reflect.Type,
	t reflect.Type,
	data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String || t.String() != "color.Color" {
		return data, nil
	}

	colorStr := data.(string)
	if colorStr == "" {
		colorStr = "black"
	}

	// Parse hex colors
	if colorStr[0] == '#' {
		return parseHexColor(colorStr)
	}

	col, ok := colornames.Map[colorStr]
	if !ok {
		return nil, fmt.Errorf("unknown color name: %s", colorStr)
	}

	return col, nil
}

func widgetDecodeHook(
	f reflect.Type,
	t reflect.Type,
	data interface{}) (interface{}, error) {

	if f.Kind() != reflect.Map || t.Name() != "Widget" {
		return data, nil
	}

	var base WidgetBase

	if err := mapstructure.Decode(data, &base); err != nil {
		return nil, err
	}

	var widget Widget

	switch base.Type {
	case "text":
		widget = &Text{}
	case "gauge":
		widget = &Gauge{}
	case "line":
		widget = &Line{}
	case "image":
		widget = &Image{}
	default:
		return nil, fmt.Errorf("unknown widget type: %s", base.Type)
	}

	if err := decodeWithHooks(data, widget, DecodeHookFunc()); err != nil {
		return nil, err
	}

	return widget, nil
}
