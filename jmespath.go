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
		return nil, parseError(expression, err)
	}

	result, err := evaluator.Evaluate(node, data)
	if err != nil {
		return nil, evaluateError(err)
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
		return nil, parseError(expression, err)
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
func (e *Expression) Search(data any) (any, error) {
	result, err := evaluator.Evaluate(e.node, data)
	if err != nil {
		return nil, evaluateError(err)
	}

	return result, nil
}

func evaluateError(err error) error {
	if errors.Is(err, evaluator.ErrInvalidType) {
		return &invalidTypeError{err.Error()}
	}

	if errors.Is(err, evaluator.ErrInvalidValue) {
		return &invalidValueError{err.Error()}
	}

	if errors.Is(err, evaluator.ErrInfinity) {
		return &infinityError{}
	}

	if errors.Is(err, evaluator.ErrNotANumber) {
		return &notANumberError{}
	}

	if err, ok := err.(*evaluator.UndefinedVariableError); ok {
		return &undefinedVariableError{err.Variable}
	}

	return &evaluationFailedError{err.Error()}
}

func parseError(expression string, err error) error {
	if err, ok := err.(*parser.InvalidFunctionArgumentError); ok {
		return &invalidTypeError{err.Error()}
	}

	if err, ok := err.(*parser.InvalidFunctionCallError); ok {
		return &invalidFunctionCallError{err.Function}
	}

	if _, ok := err.(*parser.InvalidSliceStepError); ok {
		return &invalidSliceStepError{}
	}

	if err, ok := err.(*parser.UnknownFunctionError); ok {
		return &unknownFunctionError{err.Function}
	}

	return &invalidExpressionError{expression, err.Error()}
}
