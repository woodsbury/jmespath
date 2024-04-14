package parser

import "github.com/woodsbury/jmespath/internal/lexer"

func precedence(t lexer.TokenType) int {
	switch t {
	case lexer.PipeToken:
		return 2
	case lexer.OrToken:
		return 3
	case lexer.AndToken:
		return 4
	case lexer.EqualToken,
		lexer.GreaterToken,
		lexer.GreaterOrEqualToken,
		lexer.LessToken,
		lexer.LessOrEqualToken,
		lexer.NotEqualToken:
		return 5
	case lexer.AddToken,
		lexer.SubtractToken:
		return 6
	case lexer.AsteriskToken,
		lexer.DivideToken,
		lexer.IntegerDivideToken,
		lexer.ModuloToken,
		lexer.MultiplyToken:
		return 7
	case lexer.FlattenToken:
		return 8
	case lexer.ObjectWildcardToken:
		return 9
	case lexer.FilterToken:
		return 10
	case lexer.DotToken:
		return 11
	case lexer.NotToken:
		return 12
	case lexer.ArrayWildcardToken,
		lexer.OpenSqBraceToken:
		return 13
	default:
		return 0
	}
}
