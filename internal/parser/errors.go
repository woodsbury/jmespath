package parser

import "strconv"

type InvalidFunctionArgumentError struct {
	function string
	want     string
}

func (err *InvalidFunctionArgumentError) Error() string {
	return "invalid argument to function " + strconv.Quote(err.function) + " when expecting " + err.want
}

type InvalidFunctionCallError struct {
	Function string
}

func (err *InvalidFunctionCallError) Error() string {
	return "invalid call to function " + strconv.Quote(err.Function)
}

type InvalidSliceStepError struct{}

func (err *InvalidSliceStepError) Error() string {
	return "invalid slice step value"
}

type UnknownFunctionError struct {
	Function string
}

func (err *UnknownFunctionError) Error() string {
	return "call to unknown function " + strconv.Quote(err.Function)
}

type invalidIndexError struct {
	s string
}

func (err *invalidIndexError) Error() string {
	return "invalid index " + strconv.Quote(err.s)
}

type invalidJSONLiteralError struct {
	s string
}

func (err *invalidJSONLiteralError) Error() string {
	return "invalid json literal " + strconv.Quote(err.s)
}

type invalidQuotedStringError struct {
	s string
}

func (err *invalidQuotedStringError) Error() string {
	return "invalid quoted string " + strconv.Quote(err.s)
}

type unexpectedTokenError struct {
	s string
}

func (err *unexpectedTokenError) Error() string {
	return "unexpected token " + strconv.Quote(err.s)
}
