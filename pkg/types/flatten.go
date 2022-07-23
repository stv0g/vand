package types

import (
	"reflect"
	"strings"
)

// Flatten returns all keys and corresponding values of a struct in a one-level-deep map
func Flatten(in interface{}, sep string) map[string]interface{} {
	if in == nil {
		return nil
	}
	m := map[string]interface{}{}
	appendValue(m, reflect.ValueOf(in), "", sep)
	return m
}

func appendValue(m map[string]interface{}, v reflect.Value, key string, sep string) {
	// Iterate pointers until the value
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}

		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		// Recurse to fields of structs
		for i := 0; i < v.Type().NumField(); i++ {
			sf := v.Type().Field(i) // sub-field
			sv := v.Field(i)        // sub-value

			// Skip unexported fields
			if sf.PkgPath != "" && !sf.Anonymous {
				continue
			}

			sk := sf.Name // sub-key

			tag := sf.Tag.Get("json")
			if tag != "" {
				fields := strings.SplitN(tag, ",", 2)
				sk = fields[0]
			}

			if key != "" {
				sk = key + sep + sk
			}

			appendValue(m, sv, sk, sep)
		}
		return
	}

	// Set empty key for non-struct types to the kind
	if key == "" {
		key = v.Kind().String()
	}

	m[key] = v.Interface()
}
