package lexer

import "unicode/utf8"

type Lexer struct {
	expression string
	position   int
}

func NewLexer(expression string) Lexer {
	return Lexer{
		expression: expression,
	}
}

func (l *Lexer) Next(t *Token) error {
	if l.position == len(l.expression) {
		*t = Token{
			Type: EndToken,
		}

		return nil
	}

	var r rune
	var sz int
	var err error
	for {
		r, sz, err = l.decodeRune(l.position)
		if err != nil {
			return err
		}

		if r != '\t' && r != '\n' && r != '\r' && r != ' ' {
			break
		}

		l.position += sz

		if l.position == len(l.expression) {
			*t = Token{
				Type: EndToken,
			}

			return nil
		}
	}

	start := l.position

	switch r {
	case '"':
		return l.quotedIdentifier(t, start, start+sz)
	case '$':
		return l.variable(t, start, start+sz)
	case '%':
		l.position += sz
		*t = Token{
			Type:  ModuloToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '&':
		nr, nsz, err := l.decodeRune(start + sz)
		if err == nil && nr == '&' {
			l.position += sz + nsz
			*t = Token{
				Type:  AndToken,
				Value: l.expression[start:l.position],
			}

			return nil
		}

		l.position += sz
		*t = Token{
			Type:  ExpressionToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '\'':
		return l.stringLiteral(t, start, start+sz)
	case '(':
		l.position += sz
		*t = Token{
			Type:  OpenParenToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case ')':
		l.position += sz
		*t = Token{
			Type:  CloseParenToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '*':
		l.position += sz
		*t = Token{
			Type:  AsteriskToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '+':
		l.position += sz
		*t = Token{
			Type:  AddToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case ',':
		l.position += sz
		*t = Token{
			Type:  CommaToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '-':
		nr, nsz, err := l.decodeRune(start + sz)
		if err == nil && (nr >= '0' && nr <= '9') {
			return l.numberLiteral(t, start, start+sz+nsz)
		}

		l.position += sz
		*t = Token{
			Type:  SubtractToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '.':
		nr, nsz, err := l.decodeRune(start + sz)
		if err == nil && nr == '*' {
			l.position += sz + nsz
			*t = Token{
				Type:  ObjectWildcardToken,
				Value: l.expression[start:l.position],
			}

			return nil
		}

		l.position += sz
		*t = Token{
			Type:  DotToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '/':
		nr, nsz, err := l.decodeRune(start + sz)
		if err == nil && nr == '/' {
			l.position += sz + nsz
			*t = Token{
				Type:  IntegerDivideToken,
				Value: l.expression[start:l.position],
			}

			return nil
		}

		l.position += sz
		*t = Token{
			Type:  DivideToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case ':':
		l.position += sz
		*t = Token{
			Type:  ColonToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '<':
		nr, nsz, err := l.decodeRune(start + sz)
		if err == nil && nr == '=' {
			l.position += sz + nsz
			*t = Token{
				Type:  LessOrEqualToken,
				Value: l.expression[start:l.position],
			}

			return nil
		}

		l.position += sz
		*t = Token{
			Type:  LessToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '=':
		nr, nsz, err := l.decodeRune(start + sz)
		if err == nil && nr == '=' {
			l.position += sz + nsz
			*t = Token{
				Type:  EqualToken,
				Value: l.expression[start:l.position],
			}

			return nil
		}

		l.position += sz
		*t = Token{
			Type:  AssignToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '>':
		nr, nsz, err := l.decodeRune(start + sz)
		if err == nil && nr == '=' {
			l.position += sz + nsz
			*t = Token{
				Type:  GreaterOrEqualToken,
				Value: l.expression[start:l.position],
			}

			return nil
		}

		l.position += sz
		*t = Token{
			Type:  GreaterToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '@':
		l.position += sz
		*t = Token{
			Type:  CurrentToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '[':
		nr, nsz, err := l.decodeRune(start + sz)
		if err == nil {
			if nr == '*' {
				nnr, nnsz, err := l.decodeRune(start + sz + nsz)
				if err == nil && nnr == ']' {
					l.position += sz + nsz + nnsz
					*t = Token{
						Type:  ArrayWildcardToken,
						Value: l.expression[start:l.position],
					}

					return nil
				}
			} else {
				if nr == '?' {
					l.position += sz + nsz
					*t = Token{
						Type:  FilterToken,
						Value: l.expression[start:l.position],
					}

					return nil
				}

				if nr == ']' {
					l.position += sz + nsz
					*t = Token{
						Type:  FlattenToken,
						Value: l.expression[start:l.position],
					}

					return nil
				}
			}
		}

		l.position += sz
		*t = Token{
			Type:  OpenSqBraceToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case ']':
		l.position += sz
		*t = Token{
			Type:  CloseSqBraceToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '`':
		return l.jsonLiteral(t, start, start+sz)
	case '{':
		l.position += sz
		*t = Token{
			Type:  OpenBraceToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '|':
		nr, nsz, err := l.decodeRune(start + sz)
		if err == nil && nr == '|' {
			l.position += sz + nsz
			*t = Token{
				Type:  OrToken,
				Value: l.expression[start:l.position],
			}

			return nil
		}

		l.position += sz
		*t = Token{
			Type:  PipeToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '}':
		l.position += sz
		*t = Token{
			Type:  CloseBraceToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '\u00d7':
		l.position += sz
		*t = Token{
			Type:  MultiplyToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '\u00f7':
		l.position += sz
		*t = Token{
			Type:  DivideToken,
			Value: l.expression[start:l.position],
		}

		return nil
	case '\u2212':
		l.position += sz
		*t = Token{
			Type:  SubtractToken,
			Value: l.expression[start:l.position],
		}

		return nil
	}

	if r == '!' {
		nr, nsz, err := l.decodeRune(start + sz)
		if err == nil && nr == '=' {
			l.position += sz + nsz
			*t = Token{
				Type:  NotEqualToken,
				Value: l.expression[start:l.position],
			}

			return nil
		}

		l.position += sz
		*t = Token{
			Type:  NotToken,
			Value: l.expression[start:l.position],
		}

		return nil
	}

	if r >= '0' && r <= '9' {
		return l.numberLiteral(t, start, start+sz)
	}

	if r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r == '_' {
		return l.unquotedIdentifier(t, start, start+sz)
	}

	return &unexpectedRuneError{r}
}

func (l *Lexer) decodeRune(pos int) (rune, int, error) {
	r, sz := utf8.DecodeRuneInString(l.expression[pos:])
	if sz == 0 {
		return r, sz, errUnexpectedEndOfExpression
	}

	if r == utf8.RuneError {
		return r, sz, errInvalidRune
	}

	return r, sz, nil
}

func (l *Lexer) jsonLiteral(t *Token, start, next int) error {
	for {
		r, sz, err := l.decodeRune(next)
		if err != nil {
			return err
		}

		next += sz

		if r == '`' {
			l.position = next
			*t = Token{
				Type:  JSONLiteralToken,
				Value: l.expression[start:next],
			}

			return nil
		}

		if r == '\\' {
			_, sz, err := l.decodeRune(next)
			if err != nil {
				return err
			}

			next += sz
		}
	}
}

func (l *Lexer) numberLiteral(t *Token, start, next int) error {
	for {
		r, sz, err := l.decodeRune(next)
		if err == nil && (r >= '0' && r <= '9') {
			next += sz
			continue
		}

		l.position = next
		*t = Token{
			Type:  IntegerLiteralToken,
			Value: l.expression[start:next],
		}

		return nil
	}
}

func (l *Lexer) quotedIdentifier(t *Token, start, next int) error {
	for {
		r, sz, err := l.decodeRune(next)
		if err != nil {
			return err
		}

		next += sz

		if r == '"' {
			l.position = next
			*t = Token{
				Type:  QuotedIdentifierToken,
				Value: l.expression[start:next],
			}

			return nil
		}

		if r == '\\' {
			_, sz, err := l.decodeRune(next)
			if err != nil {
				return err
			}

			next += sz
		}
	}
}

func (l *Lexer) stringLiteral(t *Token, start, next int) error {
	for {
		r, sz, err := l.decodeRune(next)
		if err != nil {
			return err
		}

		next += sz

		if r == '\'' {
			l.position = next
			*t = Token{
				Type:  StringLiteralToken,
				Value: l.expression[start:next],
			}

			return nil
		}

		if r == '\\' {
			_, sz, err := l.decodeRune(next)
			if err != nil {
				return err
			}

			next += sz
		}
	}
}

func (l *Lexer) unquotedIdentifier(t *Token, start, next int) error {
	for {
		r, sz, err := l.decodeRune(next)
		if err == nil && (r >= '0' && r <= '9' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r == '_') {
			next += sz
			continue
		}

		typ := UnquotedIdentifierToken
		switch l.expression[start:next] {
		case "in":
			typ = InToken
		case "let":
			typ = LetToken
		}

		l.position = next
		*t = Token{
			Type:  typ,
			Value: l.expression[start:next],
		}

		return nil
	}
}

func (l *Lexer) variable(t *Token, start, next int) error {
	r, sz, err := l.decodeRune(next)
	if err != nil || !(r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r == '_') {
		l.position = next
		*t = Token{
			Type:  RootToken,
			Value: l.expression[start:next],
		}

		return nil
	}

	next += sz

	for {
		r, sz, err := l.decodeRune(next)
		if err == nil && (r >= '0' && r <= '9' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r == '_') {
			next += sz
			continue
		}

		l.position = next
		*t = Token{
			Type:  VariableToken,
			Value: l.expression[start:next],
		}

		return nil
	}
}
