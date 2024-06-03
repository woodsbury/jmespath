package evaluator

import (
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/woodsbury/decimal128"
	"github.com/woodsbury/jmespath/internal/parser"
)

func (e *evaluator) arrayMaxBy(value any, node parser.Node, variables *variableScope) (any, error) {
	a, ok := value.([]any)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "array",
		}
	}

	if len(a) == 0 {
		return nil, nil
	}

	max, err := e.evaluate(node, a[0], variables)
	if err != nil {
		return nil, err
	}

	index := 0

	if strMax, ok := max.(string); ok {
		for i, v := range a[1:] {
			rv, err := e.evaluate(node, v, variables)
			if err != nil {
				return nil, err
			}

			s, ok := rv.(string)
			if !ok {
				return nil, &InvalidTypeError{
					got:  reflect.TypeOf(rv),
					want: "string",
				}
			}

			if s > strMax {
				strMax = s
				index = i + 1
			}
		}

		return a[index], nil
	}

	numMax, ok := toDecimal(max)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(max),
			want: "number",
		}
	}

	for i, v := range a[1:] {
		rv, err := e.evaluate(node, v, variables)
		if err != nil {
			return nil, err
		}

		d, ok := toDecimal(rv)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(rv),
				want: "number",
			}
		}

		if d.Cmp(numMax).Greater() {
			numMax = d
			index = i + 1
		}
	}

	return a[index], nil
}

func (e *evaluator) arrayMinBy(value any, node parser.Node, variables *variableScope) (any, error) {
	a, ok := value.([]any)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "array",
		}
	}

	if len(a) == 0 {
		return nil, nil
	}

	min, err := e.evaluate(node, a[0], variables)
	if err != nil {
		return nil, err
	}

	index := 0

	if strMin, ok := min.(string); ok {
		for i, v := range a[1:] {
			rv, err := e.evaluate(node, v, variables)
			if err != nil {
				return nil, err
			}

			s, ok := rv.(string)
			if !ok {
				return nil, &InvalidTypeError{
					got:  reflect.TypeOf(rv),
					want: "string",
				}
			}

			if s < strMin {
				strMin = s
				index = i + 1
			}
		}

		return a[index], nil
	}

	numMin, ok := toDecimal(min)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(min),
			want: "number",
		}
	}

	for i, v := range a[1:] {
		rv, err := e.evaluate(node, v, variables)
		if err != nil {
			return nil, err
		}

		d, ok := toDecimal(rv)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(rv),
				want: "number",
			}
		}

		if d.Cmp(numMin).Less() {
			numMin = d
			index = i + 1
		}
	}

	return a[index], nil
}

func (e *evaluator) filter(value any, node parser.Node, variables *variableScope) (any, error) {
	a, ok := value.([]any)
	if !ok {
		return nil, nil
	}

	r := make([]any, 0, len(a))
	for _, v := range a {
		f, err := e.evaluate(node, v, variables)
		if err != nil {
			return nil, err
		}

		if isTrue(f) {
			r = append(r, v)
		}
	}

	return r, nil
}

func (e *evaluator) filterAndProjectArray(value any, filter parser.Node, node parser.Node, variables *variableScope) (any, error) {
	a, ok := value.([]any)
	if !ok {
		return nil, nil
	}

	r := make([]any, 0, len(a))
	for _, v := range a {
		f, err := e.evaluate(filter, v, variables)
		if err != nil {
			return nil, err
		}

		if isTrue(f) {
			p, err := e.evaluate(node, v, variables)
			if err != nil {
				return nil, err
			}

			if p == nil {
				continue
			}

			r = append(r, p)
		}
	}

	return r, nil
}

func (e *evaluator) flattenAndProjectArray(value any, node parser.Node, variables *variableScope) (any, error) {
	a, ok := value.([]any)
	if !ok {
		return nil, nil
	}

	r := make([]any, 0, len(a))
	for _, v := range a {
		va, ok := v.([]any)
		if ok {
			for _, i := range va {
				p, err := e.evaluate(node, i, variables)
				if err != nil {
					return nil, err
				}

				if p == nil {
					continue
				}

				r = append(r, p)
			}

			continue
		}

		p, err := e.evaluate(node, v, variables)
		if err != nil {
			return nil, err
		}

		if p == nil {
			continue
		}

		r = append(r, p)
	}

	return r, nil
}

func (e *evaluator) mapArray(value any, node parser.Node, variables *variableScope) (any, error) {
	a, ok := value.([]any)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "array",
		}
	}

	r := make([]any, len(a))
	for i, v := range a {
		p, err := e.evaluate(node, v, variables)
		if err != nil {
			return nil, err
		}

		r[i] = p
	}

	return r, nil
}

func (e *evaluator) projectArray(value any, node parser.Node, variables *variableScope) (any, error) {
	a, ok := value.([]any)
	if !ok {
		return nil, nil
	}

	r := make([]any, 0, len(a))
	for _, v := range a {
		p, err := e.evaluate(node, v, variables)
		if err != nil {
			return nil, err
		}

		if p == nil {
			continue
		}

		r = append(r, p)
	}

	return r, nil
}

type sortByNumber struct {
	items []any
	by    []decimal128.Decimal
}

func (s sortByNumber) Len() int {
	return len(s.items)
}

func (s sortByNumber) Less(i, j int) bool {
	return decimal128.Compare(s.by[i], s.by[j]) < 0
}

func (s sortByNumber) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
	s.by[i], s.by[j] = s.by[j], s.by[i]
}

type sortByString struct {
	items []any
	by    []string
}

func (s sortByString) Len() int {
	return len(s.items)
}

func (s sortByString) Less(i, j int) bool {
	return s.by[i] < s.by[j]
}

func (s sortByString) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
	s.by[i], s.by[j] = s.by[j], s.by[i]
}

func (e *evaluator) sortArrayBy(value any, node parser.Node, variables *variableScope) (any, error) {
	a, ok := value.([]any)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(value),
			want: "array",
		}
	}

	if len(a) == 0 {
		return value, nil
	}

	first, err := e.evaluate(node, a[0], variables)
	if err != nil {
		return nil, err
	}

	if s, ok := first.(string); ok {
		by := make([]string, len(a))
		by[0] = s

		for i, v := range a[1:] {
			rv, err := e.evaluate(node, v, variables)
			if err != nil {
				return nil, err
			}

			s, ok := rv.(string)
			if !ok {
				return nil, &InvalidTypeError{
					got:  reflect.TypeOf(rv),
					want: "string",
				}
			}

			by[i+1] = s
		}

		r := sortByString{
			items: slices.Clone(a),
			by:    by,
		}

		sort.Sort(r)
		return r.items, nil
	}

	d, ok := toDecimal(first)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(first),
			want: "number",
		}
	}

	by := make([]decimal128.Decimal, len(a))
	by[0] = d

	for i, v := range a[1:] {
		rv, err := e.evaluate(node, v, variables)
		if err != nil {
			return nil, err
		}

		d, ok := toDecimal(rv)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(rv),
				want: "number",
			}
		}

		by[i+1] = d
	}

	r := sortByNumber{
		items: slices.Clone(a),
		by:    by,
	}

	sort.Sort(r)
	return r.items, nil
}

func arrayMax(v any) (any, error) {
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

	if max, ok := a[0].(string); ok {
		for _, i := range a[1:] {
			s, ok := i.(string)
			if !ok {
				return nil, &InvalidTypeError{
					got:  reflect.TypeOf(i),
					want: "string",
				}
			}

			if s > max {
				max = s
			}
		}

		return max, nil
	}

	max, ok := toDecimal(a[0])
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(a[0]),
			want: "number",
		}
	}

	for _, i := range a[1:] {
		d, ok := toDecimal(i)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(i),
				want: "number",
			}
		}

		if d.Cmp(max).Greater() {
			max = d
		}
	}

	return max, nil
}

func arrayMin(v any) (any, error) {
	a, ok := v.([]any)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(a),
			want: "array",
		}
	}

	if len(a) == 0 {
		return nil, nil
	}

	if min, ok := a[0].(string); ok {
		for _, i := range a[1:] {
			s, ok := i.(string)
			if !ok {
				return nil, &InvalidTypeError{
					got:  reflect.TypeOf(i),
					want: "string",
				}
			}

			if s < min {
				min = s
			}
		}

		return min, nil
	}

	min, ok := toDecimal(a[0])
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(a[0]),
			want: "number",
		}
	}

	for _, i := range a[1:] {
		d, ok := toDecimal(i)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(i),
				want: "number",
			}
		}

		if d.Cmp(min).Less() {
			min = d
		}
	}

	return min, nil
}

func flatten(v any) any {
	a, ok := v.([]any)
	if !ok {
		return nil
	}

	r := make([]any, 0, len(a))
	for _, v := range a {
		va, ok := v.([]any)
		if ok {
			for _, i := range va {
				if i == nil {
					continue
				}

				r = append(r, i)
			}

			continue
		}

		if v == nil {
			continue
		}

		r = append(r, v)
	}

	return r
}

func index(v any, i int) any {
	a, ok := v.([]any)
	if !ok {
		return nil
	}

	if i < 0 {
		i += len(a)
		if i < 0 {
			return nil
		}
	} else if i >= len(a) {
		return nil
	}

	return a[i]
}

func pruneArray(v any) any {
	a, ok := v.([]any)
	if !ok {
		return nil
	}

	var n bool
	var r []any
	for i, va := range a {
		if n {
			if va != nil {
				r = append(r, va)
			}

			continue
		}

		if va == nil {
			if i > 0 {
				r = append(r, a[:i]...)
			}

			n = true
		}
	}

	if n {
		return r
	}

	return a
}

func sortArray(v any) (any, error) {
	a, ok := v.([]any)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(v),
			want: "array",
		}
	}

	if len(a) == 0 {
		return v, nil
	}

	r := slices.Clone(a)

	if _, ok := a[0].(string); ok {
		valid := true
		var invalidType reflect.Type
		slices.SortFunc(r, func(a, b any) int {
			sa, ok := a.(string)
			if !ok {
				valid = false
				invalidType = reflect.TypeOf(a)
				return -1
			}

			sb, ok := b.(string)
			if !ok {
				valid = false
				invalidType = reflect.TypeOf(b)
				return 1
			}

			return strings.Compare(sa, sb)
		})

		if !valid {
			return nil, &InvalidTypeError{
				got:  invalidType,
				want: "string",
			}
		}

		return r, nil
	}

	valid := true
	var invalidType reflect.Type
	slices.SortFunc(r, func(a, b any) int {
		da, ok := toDecimal(a)
		if !ok {
			valid = false
			invalidType = reflect.TypeOf(a)
			return -1
		}

		db, ok := toDecimal(b)
		if !ok {
			valid = false
			invalidType = reflect.TypeOf(b)
			return 1
		}

		return decimal128.Compare(da, db)
	})

	if !valid {
		return nil, &InvalidTypeError{
			got:  invalidType,
			want: "number",
		}
	}

	return r, nil
}
