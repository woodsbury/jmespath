package jmespath

import (
	"testing"

	"github.com/woodsbury/decimal128"
)

func TestAllocations(t *testing.T) {
	expressions := []string{
		"object",
		"object.field",
		"let $x = object in $x.field",
		"array",
		"array[0]",
		"`\"literal\"`",
		"number",
	}

	compiled := make([]*Expression, len(expressions))
	for i, expression := range expressions {
		var err error
		compiled[i], err = Compile(expression)
		if err != nil {
			t.Fatalf("Compile(%q) = %v, want <nil>", expression, err)
		}
	}

	value := map[string]any{
		"object": map[string]any{
			"field": "value",
		},
		"array":  []any{"value"},
		"number": decimal128.FromUint32(1),
	}

	for i, expression := range compiled {
		result := testing.AllocsPerRun(1, func() {
			_, err := expression.Search(value)
			if err != nil {
				t.Fatalf("Search(%v) = %v, want <nil>", value, err)
			}
		})

		if result != 0 {
			t.Errorf("%q.Search() = %.0f allocations, want 0", expressions[i], result)
		}
	}
}
