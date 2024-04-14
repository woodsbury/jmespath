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

	dir := filepath.Join("testdata", "compliance")
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
							} else {
								switch test.Error {
								case "invalid-arity":
									if !errors.Is(err, ErrInvalidArity) {
										t.Errorf("incorrect error %v from expression %q in compliance test file %s, expected: %v", err, test.Expression, name, ErrInvalidArity)
									} else {
										pass.Add(1)
									}
								case "invalid-type":
									if !errors.Is(err, ErrInvalidType) {
										t.Errorf("incorrect error %v from expression %q in compliance test file %s, expected: %v", err, test.Expression, name, ErrInvalidType)
									} else {
										pass.Add(1)
									}
								case "invalid-value":
									if !errors.Is(err, ErrInvalidValue) {
										t.Errorf("incorrect error %v from expression %q in compliance test file %s, expected: %v", err, test.Expression, name, ErrInvalidValue)
									} else {
										pass.Add(1)
									}
								case "syntax":
									if !errors.Is(err, ErrSyntax) {
										t.Errorf("incorrect error %v from expression %q in compliance test file %s, expected %v", err, test.Expression, name, ErrSyntax)
									} else {
										pass.Add(1)
									}
								case "undefined-variable":
									if !errors.Is(err, ErrUndefinedVariable) {
										t.Errorf("incorrect error %v from expressoin %q in compliance test file %s, expected %v", err, test.Expression, name, ErrUndefinedVariable)
									} else {
										pass.Add(1)
									}
								case "unknown-function":
									if !errors.Is(err, ErrUnknownFunction) {
										t.Errorf("incorrect error %v from expression %q in compliance test file %s, expected %v", err, test.Expression, name, ErrUnknownFunction)
									} else {
										pass.Add(1)
									}
								default:
									t.Errorf("unhandled error from expression %q in compliance test file %s: %s", test.Expression, name, test.Error)
								}
							}
						} else {
							if err != nil {
								t.Errorf("unexpected error %v from expression %q in compliance test file %s", err, test.Expression, name)
							} else if !resultEqual(test.Result, result) {
								t.Errorf("incorrect result %v from expression %q in compliance test file %s, expected %v", result, test.Expression, name, test.Result)
							} else {
								pass.Add(1)
							}
						}
					}
				}
			})
		}
	})

	t.Logf("%d/%d passed", pass.Load(), total.Load())
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
	case nil:
		return y == nil
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
	}

	panic(fmt.Sprintf("unhandled type %T", x))
}
