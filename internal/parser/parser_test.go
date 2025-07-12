package parser

import "testing"

func TestParseJSONLiteral(t *testing.T) {
	t.Parallel()

	type test struct {
		value  string
		result any
	}

	tests := []test{
		{"`false`", false},
		{"`true`", true},
	}

	for _, test := range tests {
		test := test

		t.Run(test.value, func(t *testing.T) {
			node, err := parseJSONLiteral(test.value)
			if err != nil {
				t.Fatalf("parseJSONLiteral(%s) = %v, want <nil>", test.value, err)
			}

			switch node := node.(type) {
			case BoolNode:
				if result, ok := test.result.(bool); !ok || node.Value != result {
					t.Fatalf("parseJSONLiteral(%s) = %v, want %v", test.value, node.Value, test.result)
				}
			default:
				t.Fatalf("parseJSONLiteral(%s) = %v, want %v", test.value, node, test.result)
			}
		})
	}
}

func TestParseStringLiteral(t *testing.T) {
	t.Parallel()

	type test struct {
		value  string
		result string
	}

	tests := []test{
		{"''", ""},
		{"'abc'", "abc"},
		{"'\\''", "'"},
		{"'\\'", "\\"},
		{"'\\\\'", "\\"},
		{"'\\'\\\\'", "'\\"},
		{"'\\n'", "\\n"},
		{"'\\r\\n'", "\\r\\n"},
	}

	for _, test := range tests {
		test := test

		t.Run(test.value, func(t *testing.T) {
			node, err := parseStringLiteral(test.value)
			if err != nil {
				t.Fatalf("parseStringLiteral(%s) = %v, want <nil>", test.value, err)
			}

			valueNode, ok := node.(*ValueNode)
			if !ok {
				t.Fatalf("parseStringLiteral(%s) = %T, want *ValueNode", test.value, node)
			}

			stringValue, ok := valueNode.Value.(string)
			if !ok {
				t.Fatalf("parseStringLiteral(%s) = %v, want %s", test.value, valueNode.Value, test.result)
			}

			if stringValue != test.result {
				t.Fatalf("parseStringLiteral(%s) = %s, want %s", test.value, stringValue, test.result)
			}
		})
	}
}

func FuzzParser(f *testing.F) {
	f.Add("a[].b[?c == 'X'] | {x: join(', ', @)}")

	f.Fuzz(func(t *testing.T, expression string) {
		t.Parallel()

		Parse(expression)
	})
}

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()

	expression := "a[].b[?c == 'X'] | {x: join(', ', @)}"

	for b.Loop() {
		Parse(expression)
	}
}
