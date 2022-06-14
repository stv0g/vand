package widget

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/tdewolff/canvas"
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

func colorDecodeHook(
	f reflect.Type,
	t reflect.Type,
	data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String || t.String() != "color.Color" {
		return data, nil
	}

	return canvas.Red, nil
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

	fmt.Printf("unmarshaling widget: %s\n", base.Type)

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
