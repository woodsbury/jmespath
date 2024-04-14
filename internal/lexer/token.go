package lexer

type TokenType uint8

const (
	UnknownToken TokenType = iota
	EndToken

	OpenBraceToken
	CloseBraceToken
	OpenParenToken
	CloseParenToken
	OpenSqBraceToken
	CloseSqBraceToken

	AddToken
	AndToken
	ArrayWildcardToken
	AssignToken
	AsteriskToken
	ColonToken
	CommaToken
	DivideToken
	DotToken
	EqualToken
	FilterToken
	FlattenToken
	InToken
	GreaterToken
	GreaterOrEqualToken
	IntegerDivideToken
	LessToken
	LessOrEqualToken
	LetToken
	ModuloToken
	MultiplyToken
	NotToken
	NotEqualToken
	ObjectWildcardToken
	OrToken
	PipeToken
	SubtractToken

	CurrentToken
	ExpressionToken
	IntegerLiteralToken
	JSONLiteralToken
	QuotedIdentifierToken
	RootToken
	UnquotedIdentifierToken
	StringLiteralToken
	VariableToken
)

type Token struct {
	Type  TokenType
	Value string
}
