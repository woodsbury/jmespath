package lexer

import "testing"

func FuzzLexer(f *testing.F) {
	f.Add("a[].b[?c == 'X'] | {x: join(', ', @)}")

	f.Fuzz(func(t *testing.T, expression string) {
		t.Parallel()

		lex := NewLexer(expression)
		for {
			var tok Token
			if err := lex.Next(&tok); err != nil {
				return
			}

			if tok.Type == EndToken {
				return
			}
		}
	})
}

func BenchmarkLexerNext(b *testing.B) {
	b.ReportAllocs()

	expression := "a[].b[?c == 'X'] | {x: join(', ', @)}"

	for i := 0; i < b.N; i++ {
		lex := NewLexer(expression)
		for {
			var tok Token
			if err := lex.Next(&tok); err != nil {
				b.Fatal(err)
			}

			if tok.Type == EndToken {
				break
			}
		}
	}
}
