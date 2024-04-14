package evaluator

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/woodsbury/decimal128"
)

func contains(x, y any) (bool, error) {
	if x, ok := x.(string); ok {
		y, ok := y.(string)
		if !ok {
			return false, nil
		}

		return strings.Contains(x, y), nil
	}

	if x, ok := x.([]any); ok {
		for _, xi := range x {
			if equal(xi, y) {
				return true, nil
			}
		}

		return false, nil
	}

	return false, &InvalidTypeError{
		got:  reflect.TypeOf(x),
		want: "array",
	}
}

func equal(x, y any) bool {
	switch x := x.(type) {
	case nil:
		return y == nil
	case bool:
		if y, ok := y.(bool); ok {
			return x == y
		}

		return false
	case string:
		if y, ok := y.(string); ok {
			return x == y
		}

		return false
	}

	xd, ok := toDecimal(x)
	if ok {
		yd, ok := toDecimal(y)
		if !ok {
			return false
		}

		return xd.Equal(yd)
	}

	if x, ok := x.([]any); ok {
		if y, ok := y.([]any); ok {
			if len(x) != len(y) {
				return false
			}

			for i, xi := range x {
				if !equal(xi, y[i]) {
					return false
				}
			}

			return true
		}
	}

	if x, ok := x.(map[string]any); ok {
		if y, ok := y.(map[string]any); ok {
			if len(x) != len(y) {
				return false
			}

			for i, xi := range x {
				yi, ok := y[i]
				if !ok {
					return false
				}

				if !equal(xi, yi) {
					return false
				}
			}

			return true
		}
	}

	return false
}

func greater(x, y any) any {
	xd, ok := toDecimal(x)
	if !ok {
		return nil
	}

	yd, ok := toDecimal(y)
	if !ok {
		return nil
	}

	return xd.Cmp(yd).Greater()
}

func greaterOrEqual(x, y any) any {
	xd, ok := toDecimal(x)
	if !ok {
		return nil
	}

	yd, ok := toDecimal(y)
	if !ok {
		return nil
	}

	return xd.Cmp(yd).GreaterOrEqual()
}

func isTrue(v any) bool {
	switch v := v.(type) {
	case nil:
		return false
	case []any:
		return len(v) > 0
	case map[string]any:
		return len(v) > 0
	case bool:
		return v
	case float32,
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
		uint,
		decimal128.Decimal:
		return true
	case string:
		return len(v) > 0
	case json.Number:
		return len(v) > 0
	}

	return true
}

func less(x, y any) any {
	xd, ok := toDecimal(x)
	if !ok {
		return nil
	}

	yd, ok := toDecimal(y)
	if !ok {
		return nil
	}

	return xd.Cmp(yd).Less()
}

func lessOrEqual(x, y any) any {
	xd, ok := toDecimal(x)
	if !ok {
		return nil
	}

	yd, ok := toDecimal(y)
	if !ok {
		return nil
	}

	return xd.Cmp(yd).LessOrEqual()
}
