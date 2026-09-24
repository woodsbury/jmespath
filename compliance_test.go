package jmespath

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/woodsbury/decimal128"
)

func TestCompliance(t *testing.T) {
	t.Parallel()

	complianceTest(t, "compliance")
}

func TestExtra(t *testing.T) {
	t.Parallel()

	complianceTest(t, "extra")
}

func complianceTest(t *testing.T, dir string) {
	dir = filepath.Join("testdata", dir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("error reading compliance test directory: %v", err)
	}

	var total atomic.Int64
	var pass atomic.Int64

	t.Run("Files", func(t *testing.T) {
		for _, entry := range entries {
			name := entry.Name()
			if !strings.HasSuffix(name, ".json") {
				continue
			}

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				type testCase struct {
					Expression string `json:"expression"`
					Result     any    `json:"result"`
					Error      string `json:"error"`
				}

				type testCases struct {
					Given any        `json:"given"`
					Cases []testCase `json:"cases"`
				}

				f, err := os.Open(filepath.Join(dir, name))
				if err != nil {
					t.Fatalf("error reading compliance test file: %v", err)
				}

				dec := json.NewDecoder(f)
				dec.UseNumber()

				var tests []testCases
				if err := dec.Decode(&tests); err != nil {
					t.Fatalf("error decoding compliance test file: %v", err)
				}

				for _, cases := range tests {
					for _, test := range cases.Cases {
						total.Add(1)

						result, err := Search(test.Expression, cases.Given)
						if test.Error != "" {
							if err == nil {
								t.Errorf("expected error %s from expression %q in compliance test file %s", test.Error, test.Expression, name)
							} else if !errorEqual(test.Error, err) {
								t.Errorf("incorrect error %v from expression %q in compliance test file %s, expected: %s", err, test.Expression, name, test.Error)
							} else {
								pass.Add(1)
							}
						} else {
							if err != nil {
								t.Errorf("unexpected error %v from expression %q in compliance test file %s", err, test.Expression, name)
							} else if !resultEqual(test.Result, result) {
								t.Errorf("incorrect result %v from expression %q in compliance test file %s, expected %v", result, test.Expression, name, test.Result)
							} else {
								resultTypes, err := ResultTypes(test.Expression)
								if err != nil {
									t.Errorf("unexpected error %v from expression %q in compliance test file %s", err, test.Expression, name)
								} else if !resultTypeEqual(resultTypes, test.Result) {
									t.Errorf("incorrect result type from expression %q in compliance test file %s; expected %v", test.Expression, name, resultTypes)
								} else {
									pass.Add(1)
								}
							}
						}

						expr, err := Compile(test.Expression)
						if err != nil {
							if test.Error == "" {
								t.Errorf("unexpected error %v from expression %q in compliance test file %s", err, test.Expression, name)
							} else if !errorEqual(test.Error, err) {
								t.Errorf("incorrect error %v from expression %q in compliance test file %s, expected: %s", err, test.Expression, name, test.Error)
							}
						} else {
							result, err = expr.Search(cases.Given)
							if test.Error != "" {
								if err == nil {
									t.Errorf("expected error %s from expression %q in compliance test file %s", test.Error, test.Expression, name)
								} else if !errorEqual(test.Error, err) {
									t.Errorf("incorrect error %v from expression %q in compliance test file %s, expected: %s", err, test.Expression, name, test.Error)
								}
							} else {
								if err != nil {
									t.Errorf("unexpected error %v from expression %q in compliance test file %s", err, test.Expression, name)
								} else if !resultEqual(test.Result, result) {
									t.Errorf("incorrect result %v from expression %q in compliance test file %s, expected %v", result, test.Expression, name, test.Result)
								}
							}

							if resultTypes := expr.ResultTypes(); !resultTypeEqual(resultTypes, test.Result) {
								t.Errorf("incorrect result type from expression %q in compliance test file %s; expected %v", test.Expression, name, resultTypes)
							}
						}
					}
				}
			})
		}
	})

	t.Logf("%d/%d passed", pass.Load(), total.Load())
}

func errorEqual(x string, y error) bool {
	switch x {
	case "invalid-arity":
		return errors.Is(y, ErrInvalidArity)
	case "invalid-type":
		return errors.Is(y, ErrInvalidType)
	case "invalid-value":
		return errors.Is(y, ErrInvalidValue)
	case "syntax":
		return errors.Is(y, ErrSyntax)
	case "undefined-variable":
		return errors.Is(y, ErrUndefinedVariable)
	case "unknown-function":
		return errors.Is(y, ErrUnknownFunction)
	default:
		return false
	}
}

func resultEqual(x, y any) bool {
	switch x := x.(type) {
	case []any:
		y, ok := y.([]any)
		if !ok {
			return false
		}

		if len(x) != len(y) {
			return false
		}

		for i := range x {
			if !resultEqual(x[i], y[i]) {
				return false
			}
		}

		return true
	case map[string]any:
		y, ok := y.(map[string]any)
		if !ok {
			return false
		}

		if len(x) != len(y) {
			return false
		}

		for k := range x {
			if _, ok := y[k]; !ok {
				return false
			}

			if !resultEqual(x[k], y[k]) {
				return false
			}
		}

		return true
	case bool:
		y, ok := y.(bool)
		if !ok {
			return false
		}

		return x == y
	case string:
		y, ok := y.(string)
		if !ok {
			return false
		}

		return x == y
	case json.Number:
		switch y := y.(type) {
		case int64:
			x, err := x.Int64()
			if err != nil {
				panic("error parsing number: " + err.Error())
			}

			return x == y
		case json.Number:
			return x == y
		case decimal128.Decimal:
			x := decimal128.MustParse(x.String())
			return x.Equal(y)
		default:
			return false
		}
	case nil:
		return y == nil
	}

	panic(fmt.Sprintf("unhandled type %T", x))
}

func resultTypeEqual(x Types, y any) bool {
	switch y := y.(type) {
	case []any:
		return x&TypeArray != 0
	case map[string]any:
		return x&TypeObject != 0
	case bool:
		return x&TypeBoolean != 0
	case string:
		return x&TypeString != 0
	case json.Number:
		return x&TypeNumber != 0
	case nil:
		return true
	default:
		panic(fmt.Sprintf("unhandled type %T", y))
	}
}
