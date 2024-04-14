package evaluator

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/woodsbury/decimal128"
)

var (
	ErrInfinity          = errors.New("result of operation is an infinity")
	ErrInvalidType       = errors.New("invalid type")
	ErrInvalidValue      = errors.New("invalid value")
	ErrNotANumber        = errors.New("result of operation is not a number")
	ErrUndefinedVariable = errors.New("undefined variable")
)

type InvalidTypeError struct {
	got  reflect.Type
	want string
}

func (err *InvalidTypeError) Error() string {
	t := "nil"
	if err.got != nil {
		t = err.got.String()
	}

	if err.want != "" {
		return "invalid type " + t + " when expecting " + err.want
	}

	return "invalid type " + t
}

func (err *InvalidTypeError) Is(target error) bool {
	return target == ErrInvalidType
}

type UndefinedVariableError struct {
	Variable string
}

func (err *UndefinedVariableError) Error() string {
	return "undefined variable " + strconv.Quote(err.Variable)
}

func (err *UndefinedVariableError) Is(target error) bool {
	return target == ErrUndefinedVariable
}

type fromItemsKeyTypeError struct {
	key reflect.Type
}

func (err *fromItemsKeyTypeError) Error() string {
	return "array passed to from_items contains an item with a key of type " + err.key.String()
}

func (err *fromItemsKeyTypeError) Is(target error) bool {
	return target == ErrInvalidValue
}

type fromItemsLengthError struct {
	length int
}

func (err *fromItemsLengthError) Error() string {
	return "array passed to from_items contains an item of length " + strconv.Itoa(err.length)
}

func (err *fromItemsLengthError) Is(target error) bool {
	return target == ErrInvalidValue
}

type integerConversionError struct {
	num decimal128.Decimal
}

func (err *integerConversionError) Error() string {
	return "error converting value to integer: " + err.num.String()
}

func (err *integerConversionError) Is(target error) bool {
	return target == ErrInvalidValue
}

type negativeIntegerError struct {
	i int
}

func (err *negativeIntegerError) Error() string {
	return "negative integer " + strconv.Itoa(err.i) + " where positive integer required"
}

func (err *negativeIntegerError) Is(target error) bool {
	return target == ErrInvalidValue
}

type padLengthError struct {
	pad string
}

func (err *padLengthError) Error() string {
	return "padding " + strconv.Quote(err.pad) + " must have a length of 1"
}

func (err *padLengthError) Is(target error) bool {
	return target == ErrInvalidValue
}

type stringConversionError struct {
	err error
}

func (err *stringConversionError) Error() string {
	return "error converting value to string: " + err.err.Error()
}

func (err *stringConversionError) Unwrap() error {
	return err.err
}

type unexpectedOperationError struct {
	op reflect.Type
}

func (err *unexpectedOperationError) Error() string {
	var name string
	if err.op.Kind() == reflect.Pointer {
		name = err.op.Elem().Name()
	} else {
		name = err.op.Name()
	}

	return "unexpected operation " + name + " while evaluating expression"
}
