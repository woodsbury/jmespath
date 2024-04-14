package jmespath

import (
	"errors"
	"strconv"
)

var (
	// ErrEvaluationFailed indicates that the evaluation of the expression
	// failed with an error.
	ErrEvaluationFailed = errors.New("jmespath: evaluation failed")

	// ErrInvalidArity indicates that the expression called a function with an
	// incorrect number of arguments.
	ErrInvalidArity = errors.New("jmespath: invalid arity")

	// ErrInvalidType indicates that a field was used in a context where its
	// type wasn't valid.
	ErrInvalidType = errors.New("jmespath: invalid type")

	// ErrInvalidValue indicates that a field was used in a context where its
	// value wasn't valid.
	ErrInvalidValue = errors.New("jmespath: invalid value")

	// ErrNotANumber indicates that the an operation produced an infinity or
	// not-a-number result.
	ErrNotANumber = errors.New("jmespath: not a number")

	// ErrSyntax indicates that the expression contains a syntax error.
	ErrSyntax = errors.New("jmespath: syntax error")

	// ErrUndefinedVariable indicates that the expression attempted to access a
	// variable that hadn't been defined.
	ErrUndefinedVariable = errors.New("jmespath: undefined variable")

	// ErrUnknownFunction indicates that the expression attempted to invoke a
	// function that hasn't been defined.
	ErrUnknownFunction = errors.New("jmespath: unknown function")
)

type evaluationFailedError struct {
	msg string
}

func (err *evaluationFailedError) Error() string {
	return "jmespath: evaluation failed: " + err.msg
}

func (err *evaluationFailedError) Is(target error) bool {
	return target == ErrEvaluationFailed
}

type infinityError struct{}

func (err *infinityError) Error() string {
	return "jmespath: result of operation is an infinity"
}

func (err *infinityError) Is(target error) bool {
	return target == ErrNotANumber
}

type invalidFunctionCallError struct {
	function string
}

func (err *invalidFunctionCallError) Error() string {
	return "jmespath: invalid call to funcation " + strconv.Quote(err.function)
}

func (err *invalidFunctionCallError) Is(target error) bool {
	return target == ErrInvalidArity
}

type invalidExpressionError struct {
	expression string
	msg        string
}

func (err *invalidExpressionError) Error() string {
	return "jmespath: invalid expression " + strconv.Quote(err.expression) + ": " + err.msg
}

func (err *invalidExpressionError) Is(target error) bool {
	return target == ErrSyntax
}

type invalidSliceStepError struct{}

func (err *invalidSliceStepError) Error() string {
	return "jmespath: invalid slice step value"
}

func (err *invalidSliceStepError) Is(target error) bool {
	return target == ErrInvalidValue
}

type invalidTypeError struct {
	msg string
}

func (err *invalidTypeError) Error() string {
	return "jmespath: " + err.msg
}

func (err *invalidTypeError) Is(target error) bool {
	return target == ErrInvalidType
}

type invalidValueError struct {
	msg string
}

func (err *invalidValueError) Error() string {
	return "jmespath: " + err.msg
}

func (err *invalidValueError) Is(target error) bool {
	return target == ErrInvalidValue
}

type notANumberError struct{}

func (err *notANumberError) Error() string {
	return "jmespath: result of operation is not a number"
}

func (err *notANumberError) Is(target error) bool {
	return target == ErrNotANumber
}

type undefinedVariableError struct {
	variable string
}

func (err *undefinedVariableError) Error() string {
	return "jmespath: undefined variable " + strconv.Quote(err.variable)
}

func (err *undefinedVariableError) Is(target error) bool {
	return target == ErrUndefinedVariable
}

type unknownFunctionError struct {
	function string
}

func (err *unknownFunctionError) Error() string {
	return "jmespath: unknown function " + strconv.Quote(err.function)
}

func (err *unknownFunctionError) Is(target error) bool {
	return target == ErrUnknownFunction
}
