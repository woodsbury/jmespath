package evaluator

import (
	"strings"
	"unicode/utf8"

	"github.com/woodsbury/jmespath/internal/parser"
)

func isSliceNode(node parser.Node) bool {
	switch node.(type) {
	case *parser.SliceNode,
		*parser.SliceCurrentNode,
		*parser.SliceStepNode,
		*parser.SliceStepCurrentNode:
		return true
	}

	return false
}

func slice(v any, start, stop int) any {
	if a, ok := v.([]any); ok {
		l := len(a)

		if start < 0 {
			if start < -l {
				start = 0
			} else {
				start += l
			}
		} else if start >= l {
			return []any{}
		}

		if stop < 0 {
			if stop < -l {
				return []any{}
			} else {
				stop += l
			}
		} else if stop >= l {
			stop = l
		}

		if start >= stop {
			return []any{}
		}

		return a[start:stop]
	}

	if s, ok := v.(string); ok {
		l := utf8.RuneCountInString(s)

		if start < 0 {
			if start < -l {
				start = 0
			} else {
				start += l
			}
		} else if start >= l {
			return ""
		}

		if stop < 0 {
			if stop < -l {
				return ""
			} else {
				stop += l
			}
		} else if stop >= l {
			stop = l
		}

		for i := 0; i < start; i++ {
			_, sz := utf8.DecodeRuneInString(s)
			s = s[sz:]
		}

		idx := 0
		for i := start; i < stop; i++ {
			_, sz := utf8.DecodeRuneInString(s)
			idx += sz
		}

		return s[:idx]
	}

	return nil
}

func sliceStep(v any, start, stop, step int) any {
	if a, ok := v.([]any); ok {
		l := len(a)

		var n int
		if step > 0 {
			if start < 0 {
				if start < -l {
					start = 0
				} else {
					start += l
				}
			} else if start >= l {
				return []any{}
			}

			if stop < 0 {
				if stop < -l {
					return []any{}
				} else {
					stop += l
				}
			} else if stop > l {
				stop = l
			}

			if start >= stop {
				return []any{}
			}

			c := stop - start
			n = c / step
			if c%step > 0 {
				n++
			}
		} else {
			if start < 0 {
				if start < -l {
					return []any{}
				} else {
					start += l
				}
			} else if start >= l {
				start = l - 1
			}

			if stop < 0 {
				if stop < -l {
					stop = -1
				} else {
					stop += l
				}
			} else if stop >= l {
				return []any{}
			}

			if start <= stop {
				return []any{}
			}

			s := step * -1
			c := start - stop
			n = c / s
			if c%s > 0 {
				n++
			}
		}

		r := make([]any, n)
		for i, j := 0, start; i < n; i, j = i+1, j+step {
			r[i] = a[j]
		}

		return r
	}

	if s, ok := v.(string); ok {
		l := utf8.RuneCountInString(s)

		var n int
		if step > 0 {
			if start < 0 {
				if start < -l {
					start = 0
				} else {
					start += l
				}
			} else if start >= l {
				return ""
			}

			if stop < 0 {
				if stop < -l {
					return ""
				} else {
					stop += l
				}
			} else if stop > l {
				stop = l
			}

			if start >= stop {
				return ""
			}

			c := stop - start
			n = c / step
			if c%step > 0 {
				n++
			}
		} else {
			if start < 0 {
				if start < -l {
					return ""
				} else {
					start += l
				}
			} else if start >= l {
				start = l - 1
			}

			if stop < 0 {
				if stop < -l {
					stop = -1
				} else {
					stop += l
				}
			} else if stop >= l {
				return ""
			}

			if start <= stop {
				return ""
			}

			s := step * -1
			c := start - stop
			n = c / s
			if c%s > 0 {
				n++
			}
		}

		var b strings.Builder
		b.Grow(n)

		if step > 0 {
			for i := 0; i < start; i++ {
				_, sz := utf8.DecodeRuneInString(s)
				s = s[sz:]
			}

			for i := 0; i < n; i++ {
				r, sz := utf8.DecodeRuneInString(s)
				s = s[sz:]
				b.WriteRune(r)

				for j := 1; j < step; j++ {
					_, sz = utf8.DecodeRuneInString(s)
					s = s[sz:]
				}
			}
		} else {
			for i := l - 1; i > start; i-- {
				_, sz := utf8.DecodeLastRuneInString(s)
				s = s[:len(s)-sz]
			}

			for i := 0; i < n; i++ {
				r, sz := utf8.DecodeLastRuneInString(s)
				s = s[:len(s)-sz]
				b.WriteRune(r)

				for j := -1; j > step; j-- {
					_, sz = utf8.DecodeLastRuneInString(s)
					s = s[:len(s)-sz]
				}
			}
		}

		return b.String()
	}

	return nil
}
