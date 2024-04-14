package evaluator

import (
	"reflect"

	"github.com/woodsbury/jmespath/internal/parser"
)

func (e *evaluator) groupBy(value any, node parser.Node, variables *variableScope) (any, error) {
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

	r := make(map[string]any, len(a))
	for _, v := range a {
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

		if _, ok := r[s]; !ok {
			r[s] = []any{v}
		} else {
			r[s] = append(r[s].([]any), v)
		}
	}

	return r, nil
}

func (e *evaluator) projectObject(value any, node parser.Node, variables *variableScope) (any, error) {
	m, ok := value.(map[string]any)
	if !ok {
		return nil, nil
	}

	r := make([]any, 0, len(m))
	for _, v := range m {
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

func field(field string, value any) any {
	m, ok := value.(map[string]any)
	if !ok {
		return nil
	}

	return m[field]
}

func fromItems(v any) (any, error) {
	a, ok := v.([]any)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(v),
			want: "array",
		}
	}

	r := make(map[string]any, len(a))
	for _, i := range a {
		ia, ok := i.([]any)
		if !ok {
			return nil, &InvalidTypeError{
				got:  reflect.TypeOf(i),
				want: "array",
			}
		}

		if len(ia) != 2 {
			return nil, &fromItemsLengthError{
				length: len(ia),
			}
		}

		k, ok := ia[0].(string)
		if !ok {
			return nil, &fromItemsKeyTypeError{
				key: reflect.TypeOf(ia[0]),
			}
		}

		r[k] = ia[1]
	}

	return r, nil
}

func items(v any) (any, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(v),
			want: "object",
		}
	}

	r := make([]any, len(m))
	i := 0
	for k, v := range m {
		r[i] = []any{k, v}
		i++
	}

	return r, nil
}

func keys(v any) (any, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(v),
			want: "object",
		}
	}

	r := make([]any, len(m))
	i := 0
	for k := range m {
		r[i] = k
		i++
	}

	return r, nil
}

func objectValues(v any) any {
	m, ok := v.(map[string]any)
	if !ok {
		return nil
	}

	r := make([]any, len(m))
	i := 0
	for _, v := range m {
		r[i] = v
		i++
	}

	return r
}

func values(v any) (any, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return nil, &InvalidTypeError{
			got:  reflect.TypeOf(v),
			want: "object",
		}
	}

	r := make([]any, len(m))
	i := 0
	for _, v := range m {
		r[i] = v
		i++
	}

	return r, nil
}
