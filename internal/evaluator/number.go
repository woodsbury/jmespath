package evaluator

import (
	"encoding/json"
	"math"
	"reflect"

	"github.com/woodsbury/decimal128"
)

func abs(v any) (any, error) {
	if f, ok := toFloat(v); ok {
		return math.Abs(f), nil
	}

	d, ok := toDecimal(v)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(v),
			want: "number",
		}
	}

	return decimal128.Abs(d), nil
}

func add(x, y any) (any, error) {
	if xf, yf, ok := toFloatPair(x, y); ok {
		r := xf + yf

		if math.IsInf(r, 0) {
			return nil, ErrInfinity
		}

		if math.IsNaN(r) {
			return nil, ErrNotANumber
		}

		return r, nil
	}

	xd, ok := toDecimal(x)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(x),
			want: "number",
		}
	}

	yd, ok := toDecimal(y)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(y),
			want: "number",
		}
	}

	r := xd.Add(yd)

	if r.IsInf(0) {
		return nil, ErrInfinity
	}

	if r.IsNaN() {
		return nil, ErrNotANumber
	}

	return r, nil
}

func avg(v any) (any, error) {
	a, ok := v.([]any)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(v),
			want: "array",
		}
	}

	if len(a) == 0 {
		return nil, nil
	}

	var r decimal128.Decimal
	for _, v := range a {
		d, ok := toDecimal(v)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(v),
				want: "number",
			}
		}

		r = r.Add(d)
	}

	r = r.Quo(decimal128.FromInt64(int64(len(a))))

	if r.IsInf(0) {
		return nil, ErrInfinity
	}

	if r.IsNaN() {
		return nil, ErrNotANumber
	}

	return r, nil
}

func ceil(v any) (any, error) {
	if f, ok := toFloat(v); ok {
		return math.Ceil(f), nil
	}

	d, ok := toDecimal(v)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(v),
			want: "number",
		}
	}

	return decimal128.Ceil(d), nil
}

func divide(x, y any) (any, error) {
	if xf, yf, ok := toFloatPair(x, y); ok {
		r := xf / yf

		if math.IsInf(r, 0) {
			return nil, ErrInfinity
		}

		if math.IsNaN(r) {
			return nil, ErrNotANumber
		}

		return r, nil
	}

	xd, ok := toDecimal(x)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(x),
			want: "number",
		}
	}

	yd, ok := toDecimal(y)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(y),
			want: "number",
		}
	}

	r := xd.Quo(yd)

	if r.IsInf(0) {
		return nil, ErrInfinity
	}

	if r.IsNaN() {
		return nil, ErrNotANumber
	}

	return r, nil
}

func floor(v any) (any, error) {
	if f, ok := toFloat(v); ok {
		return math.Floor(f), nil
	}

	d, ok := toDecimal(v)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(v),
			want: "number",
		}
	}

	return decimal128.Floor(d), nil
}

func integerDivide(x, y any) (any, error) {
	if xf, yf, ok := toFloatPair(x, y); ok {
		r := math.Floor(xf / yf)

		if math.IsInf(r, 0) {
			return nil, ErrInfinity
		}

		if math.IsNaN(r) {
			return nil, ErrNotANumber
		}

		return r, nil
	}

	xd, ok := toDecimal(x)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(x),
			want: "number",
		}
	}

	yd, ok := toDecimal(y)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(y),
			want: "number",
		}
	}

	r, _ := xd.QuoRem(yd)

	if r.IsInf(0) {
		return nil, ErrInfinity
	}

	if r.IsNaN() {
		return nil, ErrNotANumber
	}

	return r, nil
}

func isNumber(v any) bool {
	switch v.(type) {
	case decimal128.Decimal:
		return true
	case json.Number:
		return true
	case float32:
		return true
	case float64:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case int:
		return true
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true
	case uint:
		return true
	}

	return false
}

func modulo(x, y any) (any, error) {
	if xf, yf, ok := toFloatPair(x, y); ok {
		r := math.Mod(xf, yf)

		if math.IsInf(r, 0) {
			return nil, ErrInfinity
		}

		if math.IsNaN(r) {
			return nil, ErrNotANumber
		}

		return r, nil
	}

	xd, ok := toDecimal(x)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(x),
			want: "number",
		}
	}

	yd, ok := toDecimal(y)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(y),
			want: "number",
		}
	}

	_, r := xd.QuoRem(yd)

	if r.IsInf(0) {
		return nil, ErrInfinity
	}

	if r.IsNaN() {
		return nil, ErrNotANumber
	}

	return r, nil
}

func multiply(x, y any) (any, error) {
	if xf, yf, ok := toFloatPair(x, y); ok {
		r := xf * yf

		if math.IsInf(r, 0) {
			return nil, ErrInfinity
		}

		if math.IsNaN(r) {
			return nil, ErrNotANumber
		}

		return r, nil
	}

	xd, ok := toDecimal(x)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(x),
			want: "number",
		}
	}

	yd, ok := toDecimal(y)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(y),
			want: "number",
		}
	}

	r := xd.Mul(yd)

	if r.IsInf(0) {
		return nil, ErrInfinity
	}

	if r.IsNaN() {
		return nil, ErrNotANumber
	}

	return r, nil
}

func subtract(x, y any) (any, error) {
	if xf, yf, ok := toFloatPair(x, y); ok {
		r := xf - yf

		if math.IsInf(r, 0) {
			return nil, ErrInfinity
		}

		if math.IsNaN(r) {
			return nil, ErrNotANumber
		}

		return r, nil
	}

	xd, ok := toDecimal(x)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(x),
			want: "number",
		}
	}

	yd, ok := toDecimal(y)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(y),
			want: "number",
		}
	}

	r := xd.Sub(yd)

	if r.IsInf(0) {
		return nil, ErrInfinity
	}

	if r.IsNaN() {
		return nil, ErrNotANumber
	}

	return r, nil
}

func sum(v any) (any, error) {
	a, ok := v.([]any)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(v),
			want: "array",
		}
	}

	var r decimal128.Decimal
	for _, v := range a {
		d, ok := toDecimal(v)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(v),
				want: "number",
			}
		}

		r = r.Add(d)
	}

	if r.IsInf(0) {
		return nil, ErrInfinity
	}

	if r.IsNaN() {
		return nil, ErrNotANumber
	}

	return r, nil
}

func toDecimal(v any) (decimal128.Decimal, bool) {
	switch v := v.(type) {
	case decimal128.Decimal:
		return v, true
	case json.Number:
		d, err := decimal128.Parse(v.String())
		if err != nil {
			return decimal128.Decimal{}, false
		}

		return d, true
	case float32:
		return decimal128.FromFloat32(v), true
	case float64:
		return decimal128.FromFloat64(v), true
	case int8:
		return decimal128.FromInt32(int32(v)), true
	case int16:
		return decimal128.FromInt32(int32(v)), true
	case int32:
		return decimal128.FromInt32(v), true
	case int64:
		return decimal128.FromInt64(v), true
	case int:
		return decimal128.FromInt64(int64(v)), true
	case uint8:
		return decimal128.FromUint32(uint32(v)), true
	case uint16:
		return decimal128.FromUint32(uint32(v)), true
	case uint32:
		return decimal128.FromUint32(v), true
	case uint64:
		return decimal128.FromUint64(v), true
	case uint:
		return decimal128.FromUint64(uint64(v)), true
	}

	return decimal128.Decimal{}, false
}

func toFloat(v any) (float64, bool) {
	switch v := v.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0.0, false
	}
}

func toFloatPair(x, y any) (float64, float64, bool) {
	var xf float64
	switch x := x.(type) {
	case float32:
		xf = float64(x)
	case float64:
		xf = x
	default:
		return 0.0, 0.0, false
	}

	switch y := y.(type) {
	case float32:
		return xf, float64(y), true
	case float64:
		return xf, y, true
	default:
		return 0.0, 0.0, false
	}
}

func toInt(v any) (int, bool, bool) {
	switch v := v.(type) {
	case decimal128.Decimal:
		i, ok := v.Int64()
		if !ok {
			return 0, true, false
		}

		if i > math.MaxInt || i < math.MinInt {
			return 0, true, false
		}

		return int(i), true, true
	case json.Number:
		i, err := v.Int64()
		if err != nil {
			if _, err = v.Float64(); err != nil {
				return 0, false, false
			}

			return 0, true, false
		}

		if i > math.MaxInt || i < math.MinInt {
			return 0, true, false
		}

		return int(i), true, true
	case float32:
		if v > math.MaxInt || v < math.MinInt {
			return 0, true, false
		}

		if float64(v) != math.Floor(float64(v)) {
			return 0, true, false
		}

		return int(v), true, true
	case float64:
		if v > math.MaxInt || v < math.MinInt {
			return 0, true, false
		}

		if v != math.Floor(v) {
			return 0, true, false
		}

		return int(v), true, true
	case int8:
		return int(v), true, true
	case int16:
		return int(v), true, true
	case int32:
		return int(v), true, true
	case int64:
		if v > math.MaxInt {
			return 0, true, false
		}

		return int(v), true, true
	case int:
		return v, true, true
	case uint8:
		return int(v), true, true
	case uint16:
		return int(v), true, true
	case uint32:
		return int(v), true, true
	case uint64:
		if v > math.MaxInt {
			return 0, true, false
		}

		return int(v), true, true
	case uint:
		if v > math.MaxInt {
			return 0, true, false
		}

		return int(v), true, true
	}

	return 0, false, false
}
