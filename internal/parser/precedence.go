package parser

import "github.com/woodsbury/jmespath/internal/lexer"

func precedence(t lexer.TokenType) int {
	switch t {
	case lexer.PipeToken:
		return 2
	case lexer.IfToken:
		return 3
	case lexer.OrToken:
		return 4
	case lexer.AndToken:
		return 5
	case lexer.EqualToken,
		lexer.GreaterToken,
		lexer.GreaterOrEqualToken,
		lexer.LessToken,
		lexer.LessOrEqualToken,
		lexer.NotEqualToken:
		return 6
	case lexer.AddToken,
		lexer.SubtractToken:
		return 7
	case lexer.AsteriskToken,
		lexer.DivideToken,
		lexer.IntegerDivideToken,
		lexer.ModuloToken,
		lexer.MultiplyToken:
		return 8
	case lexer.FlattenToken:
		return 9
	case lexer.ObjectWildcardToken:
		return 10
	case lexer.FilterToken:
		return 11
	case lexer.DotToken:
		return 12
	case lexer.NotToken:
		return 13
	case lexer.ArrayWildcardToken,
		lexer.OpenSqBraceToken:
		return 14
	default:
		return 0
	}
}
