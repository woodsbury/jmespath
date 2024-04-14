package evaluator

import (
	"encoding/json"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/woodsbury/decimal128"
)

func length(v any) (any, error) {
	switch v := v.(type) {
	case []any:
		return int64(len(v)), nil
	case map[string]any:
		return int64(len(v)), nil
	case string:
		return int64(utf8.RuneCountInString(v)), nil
	}

	return nil, &InvalidTypeError{
		got:  reflect.TypeOf(v),
		want: "array",
	}
}

func lower(v any) (any, error) {
	if s, ok := v.(string); ok {
		return strings.ToLower(s), nil
	}

	return nil, &InvalidTypeError{
		got:  reflect.TypeOf(v),
		want: "string",
	}
}

func reverse(v any) (any, error) {
	if s, ok := v.(string); ok {
		var b strings.Builder
		b.Grow(len(s))

		for len(s) > 0 {
			r, sz := utf8.DecodeLastRuneInString(s)
			b.WriteRune(r)
			s = s[:len(s)-sz]
		}

		return b.String(), nil
	}

	if a, ok := v.([]any); ok {
		l := len(a)
		r := make([]any, l)
		for i, j := 0, l-1; i < l; i, j = i+1, j-1 {
			r[j] = a[i]
		}

		return r, nil
	}

	return nil, &InvalidTypeError{
		got:  reflect.TypeOf(v),
		want: "array",
	}
}

func toArray(v any) any {
	if a, ok := v.([]any); ok {
		return a
	}

	return []any{v}
}

func toNumber(v any) any {
	switch v := v.(type) {
	case decimal128.Decimal,
		json.Number,
		float32,
		float64,
		int8,
		int16,
		int32,
		int64,
		int,
		uint8,
		uint16,
		uint32,
		uint64,
		uint:
		return v
	case string:
		var d decimal128.Decimal
		if err := d.UnmarshalJSON([]byte(v)); err != nil {
			return nil
		}

		return d
	}

	return nil
}

func toString(v any) (any, error) {
	if s, ok := v.(string); ok {
		return s, nil
	}

	s, err := json.Marshal(v)
	if err != nil {
		return nil, &stringConversionError{err}
	}

	return string(s), nil
}

func typeName(v any) (any, error) {
	switch v.(type) {
	case []any:
		return "array", nil
	case map[string]any:
		return "object", nil
	case bool:
		return "boolean", nil
	case decimal128.Decimal,
		json.Number,
		float32,
		float64,
		int8,
		int16,
		int32,
		int64,
		int,
		uint8,
		uint16,
		uint32,
		uint64,
		uint:
		return "number", nil
	case string:
		return "string", nil
	case nil:
		return "null", nil
	}

	return nil, &InvalidTypeError{
		got: reflect.TypeOf(v),
	}
}

func upper(v any) (any, error) {
	if s, ok := v.(string); ok {
		return strings.ToUpper(s), nil
	}

	return nil, &InvalidTypeError{
		got:  reflect.TypeOf(v),
		want: "string",
	}
}
