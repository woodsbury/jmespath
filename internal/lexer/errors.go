package lexer

import (
	"errors"
	"strconv"
)

var (
	errUnexpectedEndOfExpression = errors.New("unexpected end of expression")
)

type invalidRuneError struct {
	r rune
}

func (err *invalidRuneError) Error() string {
	return "invalid rune " + strconv.QuoteRune(err.r)
}

type unexpectedRuneError struct {
	r rune
}

func (err *unexpectedRuneError) Error() string {
	return "unexpected rune " + strconv.QuoteRune(err.r)
}
