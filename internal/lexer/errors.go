package lexer

import (
	"errors"
	"strconv"
)

var (
	errInvalidRune               = errors.New("invalid rune")
	errUnexpectedEndOfExpression = errors.New("unexpected end of expression")
)

type unexpectedRuneError struct {
	r rune
}

func (err *unexpectedRuneError) Error() string {
	return "unexpected rune " + strconv.QuoteRune(err.r)
}
