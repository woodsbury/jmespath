package jmespath

import (
	"errors"
	"strconv"

	"github.com/woodsbury/jmespath/internal/evaluator"
	"github.com/woodsbury/jmespath/internal/parser"
)

// Search evaluates expression with data and returns the result.
func Search(expression string, data any) (any, error) {
	node, err := parser.Parse(expression)
	if err != nil {
		if err, ok := err.(*parser.InvalidFunctionArgumentError); ok {
			return nil, &invalidTypeError{err.Error()}
		}

		if err, ok := err.(*parser.InvalidFunctionCallError); ok {
			return nil, &invalidFunctionCallError{err.Function}
		}

		if _, ok := err.(*parser.InvalidSliceStepError); ok {
			return nil, &invalidSliceStepError{}
		}

		if err, ok := err.(*parser.UnknownFunctionError); ok {
			return nil, &unknownFunctionError{err.Function}
		}

		return nil, &invalidExpressionError{expression, err.Error()}
	}

	result, err := evaluator.Evaluate(node, data)
	if err != nil {
		if errors.Is(err, evaluator.ErrInvalidType) {
			return nil, &invalidTypeError{err.Error()}
		}

		if errors.Is(err, evaluator.ErrInvalidValue) {
			return nil, &invalidValueError{err.Error()}
		}

		if errors.Is(err, evaluator.ErrInfinity) {
			return nil, &infinityError{}
		}

		if errors.Is(err, evaluator.ErrNotANumber) {
			return nil, &notANumberError{}
		}

		if err, ok := err.(*evaluator.UndefinedVariableError); ok {
			return nil, &undefinedVariableError{err.Variable}
		}

		return nil, &evaluationFailedError{err.Error()}
	}

	return result, nil
}

// Expression represents a compiled expression.
type Expression struct {
	node parser.Node
}

// Compile compiles expression and, if successful, returns an [Expression] that
// can be used evaluate it against data.
func Compile(expression string) (*Expression, error) {
	node, err := parser.Parse(expression)
	if err != nil {
		return nil, &invalidExpressionError{expression, err.Error()}
	}

	return &Expression{
		node: node,
	}, nil
}

// MustCompile is like [Compile] but panics if the expression cannot be
// compiled.
func MustCompile(expression string) *Expression {
	node, err := parser.Parse(expression)
	if err != nil {
		panic("jmespath.MustCompile(" + strconv.Quote(expression) + "): invalid expression")
	}

	return &Expression{
		node: node,
	}
}

// Search evaluates the compiled expression against data and returns the
// result.
func (e Expression) Search(data any) (any, error) {
	result, err := evaluator.Evaluate(e.node, data)
	if err != nil {
		if errors.Is(err, evaluator.ErrNotANumber) {
			return nil, &notANumberError{}
		}

		return nil, &evaluationFailedError{err.Error()}
	}

	return result, nil
}
