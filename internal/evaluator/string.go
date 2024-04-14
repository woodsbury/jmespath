package evaluator

import (
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

func endsWith(value, suffix any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	p, ok := suffix.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(suffix),
			want: "string",
		}
	}

	return strings.HasSuffix(s, p), nil
}

func findFirst(value, sub any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	p, ok := sub.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(sub),
			want: "string",
		}
	}

	if len(s) == 0 || len(p) == 0 {
		return nil, nil
	}

	r := strings.Index(s, p)
	if r == -1 {
		return nil, nil
	}

	r = utf8.RuneCountInString(s[:r])
	return int64(r), nil
}

func findFirstBetween(value, sub, start, finish any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	p, ok := sub.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(sub),
			want: "string",
		}
	}

	i, isNum, ok := toInt(start)
	if !ok {
		if !isNum {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(start),
				want: "number",
			}
		}

		if _, isNum, _ := toInt(finish); !isNum {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(finish),
				want: "number",
			}
		}

		d, ok := toDecimal(start)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(start),
				want: "number",
			}
		}

		return nil, &integerConversionError{
			num: d,
		}
	}

	if i < 0 {
		i = 0
	} else if i > len(s) {
		return nil, nil
	} else {
		n := 0
		for j := 0; j < i; j++ {
			_, sz := utf8.DecodeRuneInString(s[n:])
			if sz == 0 {
				return nil, nil
			}

			n += sz
		}

		i = n
	}

	j, isNum, ok := toInt(finish)
	if !ok {
		if !isNum {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(finish),
				want: "number",
			}
		}

		d, ok := toDecimal(finish)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(finish),
				want: "number",
			}
		}

		return nil, &integerConversionError{
			num: d,
		}
	}

	if j < 0 {
		return nil, nil
	} else if j > len(s) {
		j = len(s)
	} else {
		n := 0
		for k := 0; k < j; k++ {
			_, sz := utf8.DecodeRuneInString(s[n:])
			if sz == 0 {
				return nil, nil
			}

			n += sz
		}

		j = n
	}

	r := strings.Index(s[i:j], p)
	if r == -1 {
		return nil, nil
	}

	r = utf8.RuneCountInString(s[:r+i])
	return int64(r), nil
}

func findFirstFrom(value, sub, start any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	p, ok := sub.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(sub),
			want: "string",
		}
	}

	i, isNum, ok := toInt(start)
	if !ok {
		if !isNum {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(start),
				want: "number",
			}
		}

		d, ok := toDecimal(start)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(start),
				want: "number",
			}
		}

		return nil, &integerConversionError{
			num: d,
		}
	}

	if i < 0 {
		i = 0
	} else if i > len(s) {
		return nil, nil
	} else {
		n := 0
		for j := 0; j < i; j++ {
			_, sz := utf8.DecodeRuneInString(s[n:])
			if sz == 0 {
				return nil, nil
			}

			n += sz
		}

		i = n
	}

	r := strings.Index(s[i:], p)
	if r == -1 {
		return nil, nil
	}

	r = utf8.RuneCountInString(s[:r+i])
	return int64(r), nil
}

func findLast(value, sub any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	p, ok := sub.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(sub),
			want: "string",
		}
	}

	if len(s) == 0 || len(p) == 0 {
		return nil, nil
	}

	r := strings.LastIndex(s, p)
	if r == -1 {
		return nil, nil
	}

	r = utf8.RuneCountInString(s[:r])
	return int64(r), nil
}

func findLastBetween(value, sub, start, finish any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	p, ok := sub.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(sub),
			want: "string",
		}
	}

	i, isNum, ok := toInt(start)
	if !ok {
		if !isNum {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(start),
				want: "number",
			}
		}

		if _, isNum, _ := toInt(finish); !isNum {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(finish),
				want: "number",
			}
		}

		d, ok := toDecimal(start)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(start),
				want: "number",
			}
		}

		return nil, &integerConversionError{
			num: d,
		}
	}

	if i < 0 {
		i = 0
	} else if i > len(s) {
		return nil, nil
	} else {
		n := 0
		for j := 0; j < i; j++ {
			_, sz := utf8.DecodeRuneInString(s[n:])
			if sz == 0 {
				return nil, nil
			}

			n += sz
		}

		i = n
	}

	j, isNum, ok := toInt(finish)
	if !ok {
		if !isNum {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(finish),
				want: "number",
			}
		}

		d, ok := toDecimal(finish)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(finish),
				want: "number",
			}
		}

		return nil, &integerConversionError{
			num: d,
		}
	}

	if j < 0 {
		return nil, nil
	} else if j > len(s) {
		j = len(s)
	} else {
		n := 0
		for k := 0; k < j; k++ {
			_, sz := utf8.DecodeRuneInString(s[n:])
			if sz == 0 {
				return nil, nil
			}

			n += sz
		}

		j = n
	}

	r := strings.LastIndex(s[i:j], p)
	if r == -1 {
		return nil, nil
	}

	r = utf8.RuneCountInString(s[:r+i])
	return int64(r), nil
}

func findLastFrom(value, sub, start any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	p, ok := sub.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(sub),
			want: "string",
		}
	}

	i, isNum, ok := toInt(start)
	if !ok {
		if !isNum {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(start),
				want: "number",
			}
		}

		d, ok := toDecimal(start)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(start),
				want: "number",
			}
		}

		return nil, &integerConversionError{
			num: d,
		}
	}

	if i < 0 {
		i = 0
	} else if i > len(s) {
		return nil, nil
	} else {
		n := 0
		for j := 0; j < i; j++ {
			_, sz := utf8.DecodeRuneInString(s[n:])
			if sz == 0 {
				return nil, nil
			}

			n += sz
		}

		i = n
	}

	r := strings.LastIndex(s[i:], p)
	if r == -1 {
		return nil, nil
	}

	r = utf8.RuneCountInString(s[:r+i])
	return int64(r), nil
}

func join(sep, value any) (any, error) {
	a, ok := value.([]any)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "array",
		}
	}

	s, ok := sep.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(sep),
			want: "string",
		}
	}

	if len(a) == 0 {
		return "", nil
	}

	e, ok := a[0].(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(a[0]),
			want: "string",
		}
	}

	var b strings.Builder
	b.WriteString(e)

	for _, i := range a[1:] {
		e, ok := i.(string)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(i),
				want: "string",
			}
		}

		b.WriteString(s)
		b.WriteString(e)
	}

	return b.String(), nil
}

func padLeft(value, width, pad any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	p, ok := pad.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(pad),
			want: "string",
		}
	}

	w, isNum, ok := toInt(width)
	if !ok {
		if !isNum {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(width),
				want: "number",
			}
		}

		d, ok := toDecimal(width)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(width),
				want: "number",
			}
		}

		return nil, &integerConversionError{
			num: d,
		}
	}

	if w < 0 {
		return nil, &negativeIntegerError{
			i: w,
		}
	}

	if len(p) != 1 {
		return nil, &padLengthError{
			pad: p,
		}
	}

	n := w - len(s)
	if n <= 0 {
		return value, nil
	}

	var b strings.Builder

	for n > 0 {
		b.WriteString(p)
		n--
	}

	b.WriteString(s)
	return b.String(), nil
}

func padRight(value, width, pad any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	p, ok := pad.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(pad),
			want: "string",
		}
	}

	w, isNum, ok := toInt(width)
	if !ok {
		if !isNum {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(width),
				want: "number",
			}
		}

		d, ok := toDecimal(width)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(width),
				want: "number",
			}
		}

		return nil, &integerConversionError{
			num: d,
		}
	}

	if w < 0 {
		return nil, &negativeIntegerError{
			i: w,
		}
	}

	if len(p) != 1 {
		return nil, &padLengthError{
			pad: p,
		}
	}

	n := w - len(s)
	if n <= 0 {
		return value, nil
	}

	var b strings.Builder
	b.WriteString(s)

	for n > 0 {
		b.WriteString(p)
		n--
	}

	return b.String(), nil
}

func padSpaceLeft(value, width any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	w, isNum, ok := toInt(width)
	if !ok {
		if !isNum {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(width),
				want: "number",
			}
		}

		d, ok := toDecimal(width)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(width),
				want: "number",
			}
		}

		return nil, &integerConversionError{
			num: d,
		}
	}

	if w < 0 {
		return nil, &negativeIntegerError{
			i: w,
		}
	}

	n := w - len(s)
	if n <= 0 {
		return value, nil
	}

	var b strings.Builder

	for n > 0 {
		b.WriteByte(' ')
		n--
	}

	b.WriteString(s)
	return b.String(), nil
}

func padSpaceRight(value, width any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	w, isNum, ok := toInt(width)
	if !ok {
		if !isNum {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(width),
				want: "number",
			}
		}

		d, ok := toDecimal(width)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(width),
				want: "number",
			}
		}

		return nil, &integerConversionError{
			num: d,
		}
	}

	if w < 0 {
		return nil, &negativeIntegerError{
			i: w,
		}
	}

	n := w - len(s)
	if n <= 0 {
		return value, nil
	}

	var b strings.Builder
	b.WriteString(s)

	for n > 0 {
		b.WriteByte(' ')
		n--
	}

	return b.String(), nil
}

func replace(value, old, new any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	po, ok := old.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(old),
			want: "string",
		}
	}

	pn, ok := new.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(new),
			want: "string",
		}
	}

	return strings.ReplaceAll(s, po, pn), nil
}

func replaceCount(value, old, new, count any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	po, ok := old.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(old),
			want: "string",
		}
	}

	pn, ok := new.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(new),
			want: "string",
		}
	}

	n, isNum, ok := toInt(count)
	if !ok {
		if !isNum {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(count),
				want: "number",
			}
		}

		d, ok := toDecimal(count)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(count),
				want: "number",
			}
		}

		return nil, &integerConversionError{
			num: d,
		}
	}

	return strings.Replace(s, po, pn, n), nil
}

func split(value, sep any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	p, ok := sep.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(sep),
			want: "string",
		}
	}

	if len(s) == 0 {
		return []any{}, nil
	}

	if len(p) == 0 {
		n := utf8.RuneCountInString(s) - 1
		r := make([]any, n+1)

		i := 0
		for i < n {
			_, l := utf8.DecodeRuneInString(s)
			r[i] = s[:l]
			s = s[l:]
			i++
		}

		r[i] = s
		return r[:i+1], nil
	}

	n := strings.Count(s, p)
	r := make([]any, n+1)

	i := 0
	for i < n {
		j := strings.Index(s, p)
		if j < 0 {
			break
		}

		r[i] = s[:j]
		s = s[j+len(p):]
		i++
	}

	r[i] = s
	return r[:i+1], nil
}

func splitCount(value, sep, count any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	p, ok := sep.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(sep),
			want: "string",
		}
	}

	n, isNum, ok := toInt(count)
	if !ok {
		if !isNum {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(count),
				want: "number",
			}
		}

		d, ok := toDecimal(count)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(count),
				want: "number",
			}
		}

		return nil, &integerConversionError{
			num: d,
		}
	}

	if n < 0 {
		return nil, &negativeIntegerError{
			i: n,
		}
	}

	if n == 0 {
		return []any{s}, nil
	}

	if len(s) == 0 {
		return []any{}, nil
	}

	if len(p) == 0 {
		r := make([]any, n+1)

		i := 0
		for i < n {
			_, l := utf8.DecodeRuneInString(s)
			r[i] = s[:l]
			s = s[l:]
			i++
		}

		r[i] = s
		return r[:i+1], nil
	}

	r := make([]any, n+1)

	i := 0
	for i < n {
		j := strings.Index(s, p)
		if j < 0 {
			break
		}

		r[i] = s[:j]
		s = s[j+len(p):]
		i++
	}

	r[i] = s
	return r[:i+1], nil
}

func startsWith(value, prefix any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	p, ok := prefix.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(prefix),
			want: "string",
		}
	}

	return strings.HasPrefix(s, p), nil
}

func trim(value, cut any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	p, ok := cut.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(cut),
			want: "string",
		}
	}

	if len(p) == 0 {
		return strings.TrimSpace(s), nil
	}

	return strings.Trim(s, p), nil
}

func trimLeft(value, cut any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	p, ok := cut.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(cut),
			want: "string",
		}
	}

	if len(p) == 0 {
		return strings.TrimLeftFunc(s, unicode.IsSpace), nil
	}

	return strings.TrimLeft(s, p), nil
}

func trimRight(value, cut any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	p, ok := cut.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(cut),
			want: "string",
		}
	}

	if len(p) == 0 {
		return strings.TrimRightFunc(s, unicode.IsSpace), nil
	}

	return strings.TrimRight(s, p), nil
}

func trimSpace(value any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	return strings.TrimSpace(s), nil
}

func trimSpaceLeft(value any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	return strings.TrimLeftFunc(s, unicode.IsSpace), nil
}

func trimSpaceRight(value any) (any, error) {
	s, ok := value.(string)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "string",
		}
	}

	return strings.TrimRightFunc(s, unicode.IsSpace), nil
}
