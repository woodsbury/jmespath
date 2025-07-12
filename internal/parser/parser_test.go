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
			result, err := parseJSONLiteral(test.value)
			if err != nil || result != test.result {
				t.Fatalf("parseJSONLiteral(%s) = (%v, %v), want (%v, <nil>)", test.value, result, err, test.result)
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
			result, err := parseStringLiteral(test.value)
			if err != nil || result != test.result {
				t.Fatalf("parseStringLiteral(%s) = (%s, %v), want (%s, <nil>)", test.value, result, err, test.result)
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
