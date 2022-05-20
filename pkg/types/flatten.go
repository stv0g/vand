package types

import (
	"reflect"
	"strconv"
	"strings"
)

const (
	// ScopeDelimiter is the character used for joining multi-level strings
	// make configurable with a tag?
	ScopeDelimiter = "."
)

var enumerables = map[reflect.Kind]bool{reflect.Slice: true, reflect.Array: true}

// Flatten returns all keys and corresponding values of a struct in a one-level-deep map
func Flatten(in interface{}) map[string]interface{} {
	if in == nil {
		return nil
	}
	m := map[string]interface{}{}
	appendValue(m, reflect.ValueOf(in), "")
	return m
}

func appendValue(m map[string]interface{}, v reflect.Value, key string) {
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
				sk = key + ScopeDelimiter + sk
			}

			appendValue(m, sv, sk)
		}
		return
	}

	// Set empty key for non-struct types to the kind
	if key == "" {
		key = v.Kind().String()
	}

	// make 1 vs. 0 indexing configurable with tags?
	if enumerables[v.Kind()] {
		for i := 0; i < v.Len(); i++ {
			sk := key + ScopeDelimiter + strconv.Itoa(i)
			m[sk] = v.Index(i).Interface()
		}
		return
	}

	m[key] = v.Interface()
}
