package parser

import (
	"encoding/json"
	"errors"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/woodsbury/jmespath/internal/lexer"
)

func Parse(expression string) (Node, error) {
	p := parser{
		lex: lexer.NewLexer(expression),
	}

	if err := p.lex.Next(&p.curr); err != nil {
		return nil, err
	}

	if err := p.lex.Next(&p.next); err != nil {
		return nil, err
	}

	return p.parse()
}

type parser struct {
	lex  lexer.Lexer
	curr lexer.Token
	next lexer.Token
}

func (p *parser) advance() error {
	p.curr = p.next
	return p.lex.Next(&p.next)
}

func (p *parser) advance2() error {
	if err := p.lex.Next(&p.curr); err != nil {
		return err
	}

	return p.lex.Next(&p.next)
}

func (p *parser) expression(prec int) (Node, error) {
	node, err := p.primaryExpression()
	if err != nil {
		return nil, err
	}

	newPrec := precedence(p.curr.Type)
	for newPrec > prec {
		switch p.curr.Type {
		case lexer.AddToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.expression(newPrec)
			if err != nil {
				return nil, err
			}

			node = &AddNode{
				Left:  node,
				Right: right,
			}
		case lexer.AndToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.expression(newPrec)
			if err != nil {
				return nil, err
			}

			node = &AndNode{
				Left:  node,
				Right: right,
			}
		case lexer.ArrayWildcardToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.projection(precedence(lexer.ObjectWildcardToken))
			if err != nil {
				return nil, err
			}

			if right == nil {
				node = &PruneArrayNode{
					Child: node,
				}
			} else {
				node = &ProjectArrayNode{
					Left:  node,
					Right: right,
				}
			}
		case lexer.AsteriskToken,
			lexer.MultiplyToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.expression(newPrec)
			if err != nil {
				return nil, err
			}

			node = &MultiplyNode{
				Left:  node,
				Right: right,
			}
		case lexer.DivideToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.expression(newPrec)
			if err != nil {
				return nil, err
			}

			node = &DivideNode{
				Left:  node,
				Right: right,
			}
		case lexer.DotToken:
			switch p.next.Type {
			case lexer.ArrayWildcardToken:
				if err := p.advance2(); err != nil {
					return nil, err
				}

				node = &SelectArraySingleNode{
					Child: node,
					Field: ObjectValuesCurrentNode{},
				}
			case lexer.OpenBraceToken:
				if err := p.advance2(); err != nil {
					return nil, err
				}

				node, err = p.selectObject(node)
				if err != nil {
					return nil, err
				}
			case lexer.OpenSqBraceToken:
				if err := p.advance2(); err != nil {
					return nil, err
				}

				node, err = p.selectArray(node)
				if err != nil {
					return nil, err
				}
			case lexer.QuotedIdentifierToken:
				if err := p.advance(); err != nil {
					return nil, err
				}

				if isProjectNode(node) {
					right, err := p.expression(newPrec)
					if err != nil {
						return nil, err
					}

					node = &ProjectArrayNode{
						Left:  node,
						Right: right,
					}
				} else if precedence(p.next.Type) > newPrec {
					right, err := p.expression(newPrec)
					if err != nil {
						return nil, err
					}

					node = &PipeNode{
						Left:  node,
						Right: right,
					}
				} else {
					right, err := parseQuotedIdentifier(p.curr.Value)
					if err != nil {
						return nil, err
					}

					node = &PipeFieldNode{
						Left:  node,
						Right: right,
					}

					if err := p.advance(); err != nil {
						return nil, err
					}
				}
			case lexer.UnquotedIdentifierToken:
				if err := p.advance(); err != nil {
					return nil, err
				}

				if isProjectNode(node) {
					right, err := p.expression(newPrec)
					if err != nil {
						return nil, err
					}

					node = &ProjectArrayNode{
						Left:  node,
						Right: right,
					}
				} else if p.next.Type == lexer.OpenParenToken || precedence(p.next.Type) > newPrec {
					right, err := p.expression(newPrec)
					if err != nil {
						return nil, err
					}

					node = &PipeNode{
						Left:  node,
						Right: right,
					}
				} else {
					node = &PipeFieldNode{
						Left:  node,
						Right: p.curr.Value,
					}

					if err := p.advance(); err != nil {
						return nil, err
					}
				}
			default:
				return nil, &unexpectedTokenError{p.curr.Value}
			}
		case lexer.EqualToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.expression(newPrec)
			if err != nil {
				return nil, err
			}

			node = &EqualNode{
				Left:  node,
				Right: right,
			}
		case lexer.FilterToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			filter, err := p.filter()
			if err != nil {
				return nil, err
			}

			right, err := p.projection(newPrec)
			if err != nil {
				return nil, err
			}

			if right == nil {
				node = &FilterNode{
					Child:  node,
					Filter: filter,
				}
			} else {
				node = &FilterAndProjectNode{
					Left:   node,
					Filter: filter,
					Right:  right,
				}
			}
		case lexer.FlattenToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.projection(newPrec)
			if err != nil {
				return nil, err
			}

			if right == nil {
				node = &FlattenNode{
					Child: node,
				}
			} else {
				node = &FlattenAndProjectNode{
					Left:  node,
					Right: right,
				}
			}
		case lexer.GreaterToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.expression(newPrec)
			if err != nil {
				return nil, err
			}

			node = &GreaterNode{
				Left:  node,
				Right: right,
			}
		case lexer.GreaterOrEqualToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.expression(newPrec)
			if err != nil {
				return nil, err
			}

			node = &GreaterOrEqualNode{
				Left:  node,
				Right: right,
			}
		case lexer.IntegerDivideToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.expression(newPrec)
			if err != nil {
				return nil, err
			}

			node = &IntegerDivideNode{
				Left:  node,
				Right: right,
			}
		case lexer.LessToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.expression(newPrec)
			if err != nil {
				return nil, err
			}

			node = &LessNode{
				Left:  node,
				Right: right,
			}
		case lexer.LessOrEqualToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.expression(newPrec)
			if err != nil {
				return nil, err
			}

			node = &LessOrEqualNode{
				Left:  node,
				Right: right,
			}
		case lexer.ModuloToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.expression(newPrec)
			if err != nil {
				return nil, err
			}

			node = &ModuloNode{
				Left:  node,
				Right: right,
			}
		case lexer.NotEqualToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.expression(newPrec)
			if err != nil {
				return nil, err
			}

			node = &NotEqualNode{
				Left:  node,
				Right: right,
			}
		case lexer.ObjectWildcardToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.projection(newPrec)
			if err != nil {
				return nil, err
			}

			if right == nil {
				node = &ObjectValuesNode{
					Child: node,
				}
			} else {
				node = &ProjectObjectNode{
					Left:  node,
					Right: right,
				}
			}
		case lexer.OpenSqBraceToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			var project bool
			node, project, err = p.index(node)
			if err != nil {
				return nil, err
			}

			if project {
				right, err := p.projection(newPrec)
				if err != nil {
					return nil, err
				}

				if right != nil {
					node = &ProjectArrayNode{
						Left:  node,
						Right: right,
					}
				}
			}
		case lexer.OrToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.expression(newPrec)
			if err != nil {
				return nil, err
			}

			node = &OrNode{
				Left:  node,
				Right: right,
			}
		case lexer.PipeToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.expression(newPrec)
			if err != nil {
				return nil, err
			}

			node = &PipeNode{
				Left:  node,
				Right: right,
			}
		case lexer.SubtractToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			right, err := p.expression(newPrec)
			if err != nil {
				return nil, err
			}

			node = &SubtractNode{
				Left:  node,
				Right: right,
			}
		default:
			return node, nil
		}

		newPrec = precedence(p.curr.Type)
	}

	return node, nil
}

func (p *parser) filter() (Node, error) {
	node, err := p.expression(1)
	if err != nil {
		return nil, err
	}

	if p.curr.Type != lexer.CloseSqBraceToken {
		return nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, err
	}

	return node, nil
}

func (p *parser) function() (Node, error) {
	name := p.curr.Value

	if err := p.advance2(); err != nil {
		return nil, err
	}

	switch name {
	case "abs":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &AbsNode{
			Argument: arg,
		}, nil
	case "avg":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &AvgNode{
			Argument: arg,
		}, nil
	case "ceil":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &CeilNode{
			Argument: arg,
		}, nil
	case "contains":
		arg1, arg2, err := p.function2Arg(name)
		if err != nil {
			return nil, err
		}

		return &ContainsNode{
			Arguments: [2]Node{arg1, arg2},
		}, nil
	case "ends_with":
		arg1, arg2, err := p.function2Arg(name)
		if err != nil {
			return nil, err
		}

		return &EndsWithNode{
			Arguments: [2]Node{arg1, arg2},
		}, nil
	case "find_first":
		arg1, arg2, arg3, arg4, err := p.function2To4Arg(name)
		if err != nil {
			return nil, err
		}

		if arg3 == nil {
			return &FindFirstNode{
				Arguments: [2]Node{arg1, arg2},
			}, nil
		}

		if arg4 == nil {
			return &FindFirstFromNode{
				Arguments: [3]Node{arg1, arg2, arg3},
			}, nil
		}

		return &FindFirstBetweenNode{
			Arguments: [4]Node{arg1, arg2, arg3, arg4},
		}, nil
	case "find_last":
		arg1, arg2, arg3, arg4, err := p.function2To4Arg(name)
		if err != nil {
			return nil, err
		}

		if arg3 == nil {
			return &FindLastNode{
				Arguments: [2]Node{arg1, arg2},
			}, nil
		}

		if arg4 == nil {
			return &FindLastFromNode{
				Arguments: [3]Node{arg1, arg2, arg3},
			}, nil
		}

		return &FindLastBetweenNode{
			Arguments: [4]Node{arg1, arg2, arg3, arg4},
		}, nil
	case "floor":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &FloorNode{
			Argument: arg,
		}, nil
	case "from_items":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &FromItemsNode{
			Argument: arg,
		}, nil
	case "group_by":
		arg1, arg2, err := p.function2ExpArg(name)
		if err != nil {
			return nil, err
		}

		return &GroupByNode{
			Arguments: [2]Node{arg1, arg2},
		}, nil
	case "items":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &ItemsNode{
			Argument: arg,
		}, nil
	case "join":
		arg1, arg2, err := p.function2Arg(name)
		if err != nil {
			return nil, err
		}

		return &JoinNode{
			Arguments: [2]Node{arg1, arg2},
		}, nil
	case "keys":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &KeysNode{
			Argument: arg,
		}, nil
	case "length":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &LengthNode{
			Argument: arg,
		}, nil
	case "lower":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &LowerNode{
			Argument: arg,
		}, nil
	case "map":
		arg1, arg2, err := p.function2MapArg(name)
		if err != nil {
			return nil, err
		}

		return &MapNode{
			Arguments: [2]Node{arg1, arg2},
		}, nil
	case "max":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &MaxNode{
			Argument: arg,
		}, nil
	case "max_by":
		arg1, arg2, err := p.function2ExpArg(name)
		if err != nil {
			return nil, err
		}

		return &MaxByNode{
			Arguments: [2]Node{arg1, arg2},
		}, nil
	case "merge":
		args, err := p.functionVarArg(name)
		if err != nil {
			return nil, err
		}

		return &MergeNode{
			Arguments: args,
		}, nil
	case "min":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &MinNode{
			Argument: arg,
		}, nil
	case "min_by":
		arg1, arg2, err := p.function2ExpArg(name)
		if err != nil {
			return nil, err
		}

		return &MinByNode{
			Arguments: [2]Node{arg1, arg2},
		}, nil
	case "not_null":
		return p.functionNotNull()
	case "pad_left":
		arg1, arg2, arg3, err := p.function2To3Arg(name)
		if err != nil {
			return nil, err
		}

		if arg3 == nil {
			return &PadSpaceLeftNode{
				Arguments: [2]Node{arg1, arg2},
			}, nil
		}

		return &PadLeftNode{
			Arguments: [3]Node{arg1, arg2, arg3},
		}, nil
	case "pad_right":
		arg1, arg2, arg3, err := p.function2To3Arg(name)
		if err != nil {
			return nil, err
		}

		if arg3 == nil {
			return &PadSpaceRightNode{
				Arguments: [2]Node{arg1, arg2},
			}, nil
		}

		return &PadRightNode{
			Arguments: [3]Node{arg1, arg2, arg3},
		}, nil
	case "replace":
		arg1, arg2, arg3, arg4, err := p.function3To4Arg(name)
		if err != nil {
			return nil, err
		}

		if arg4 == nil {
			return &ReplaceNode{
				Arguments: [3]Node{arg1, arg2, arg3},
			}, nil
		}

		return &ReplaceCountNode{
			Arguments: [4]Node{arg1, arg2, arg3, arg4},
		}, nil
	case "reverse":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &ReverseNode{
			Argument: arg,
		}, nil
	case "sort":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &SortNode{
			Argument: arg,
		}, nil
	case "sort_by":
		arg1, arg2, err := p.function2ExpArg(name)
		if err != nil {
			return nil, err
		}

		return &SortByNode{
			Arguments: [2]Node{arg1, arg2},
		}, nil
	case "split":
		arg1, arg2, arg3, err := p.function2To3Arg(name)
		if err != nil {
			return nil, err
		}

		if arg3 == nil {
			return &SplitNode{
				Arguments: [2]Node{arg1, arg2},
			}, nil
		}

		return &SplitCountNode{
			Arguments: [3]Node{arg1, arg2, arg3},
		}, nil
	case "starts_with":
		arg1, arg2, err := p.function2Arg(name)
		if err != nil {
			return nil, err
		}

		return &StartsWithNode{
			Arguments: [2]Node{arg1, arg2},
		}, nil
	case "sum":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &SumNode{
			Argument: arg,
		}, nil
	case "to_array":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &ToArrayNode{
			Argument: arg,
		}, nil
	case "to_number":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &ToNumberNode{
			Argument: arg,
		}, nil
	case "to_string":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &ToStringNode{
			Argument: arg,
		}, nil
	case "trim":
		arg1, arg2, err := p.function1To2Arg(name)
		if err != nil {
			return nil, err
		}

		if arg2 == nil {
			return &TrimSpaceNode{
				Argument: arg1,
			}, nil
		}

		return &TrimNode{
			Arguments: [2]Node{arg1, arg2},
		}, nil
	case "trim_left":
		arg1, arg2, err := p.function1To2Arg(name)
		if err != nil {
			return nil, err
		}

		if arg2 == nil {
			return &TrimSpaceLeftNode{
				Argument: arg1,
			}, nil
		}

		return &TrimLeftNode{
			Arguments: [2]Node{arg1, arg2},
		}, nil
	case "trim_right":
		arg1, arg2, err := p.function1To2Arg(name)
		if err != nil {
			return nil, err
		}

		if arg2 == nil {
			return &TrimSpaceRightNode{
				Argument: arg1,
			}, nil
		}

		return &TrimRightNode{
			Arguments: [2]Node{arg1, arg2},
		}, nil
	case "type":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &TypeNode{
			Argument: arg,
		}, nil
	case "upper":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &UpperNode{
			Argument: arg,
		}, nil
	case "values":
		arg, err := p.function1Arg(name)
		if err != nil {
			return nil, err
		}

		return &ValuesNode{
			Argument: arg,
		}, nil
	case "zip":
		args, err := p.functionVarArg(name)
		if err != nil {
			return nil, err
		}

		return &ZipNode{
			Arguments: args,
		}, nil
	}

	return nil, &UnknownFunctionError{name}
}

func (p *parser) function1Arg(name string) (Node, error) {
	if p.curr.Type == lexer.CloseParenToken {
		return nil, &InvalidFunctionCallError{name}
	}

	arg, err := p.expression(1)
	if err != nil {
		return nil, err
	}

	if p.curr.Type == lexer.CommaToken {
		return nil, &InvalidFunctionCallError{name}
	}

	if p.curr.Type != lexer.CloseParenToken {
		return nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, err
	}

	return arg, nil
}

func (p *parser) function1To2Arg(name string) (Node, Node, error) {
	if p.curr.Type == lexer.CloseParenToken {
		return nil, nil, &InvalidFunctionCallError{name}
	}

	arg1, err := p.expression(1)
	if err != nil {
		return nil, nil, err
	}

	if p.curr.Type == lexer.CloseParenToken {
		if err := p.advance(); err != nil {
			return nil, nil, err
		}

		return arg1, nil, nil
	}

	if p.curr.Type != lexer.CommaToken {
		return nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, err
	}

	arg2, err := p.expression(1)
	if err != nil {
		return nil, nil, err
	}

	if p.curr.Type == lexer.CommaToken {
		return nil, nil, &InvalidFunctionCallError{name}
	}

	if p.curr.Type != lexer.CloseParenToken {
		return nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, err
	}

	return arg1, arg2, nil
}

func (p *parser) function2Arg(name string) (Node, Node, error) {
	if p.curr.Type == lexer.CloseParenToken {
		return nil, nil, &InvalidFunctionCallError{name}
	}

	arg1, err := p.expression(1)
	if err != nil {
		return nil, nil, err
	}

	if p.curr.Type == lexer.CloseParenToken {
		return nil, nil, &InvalidFunctionCallError{name}
	}

	if p.curr.Type != lexer.CommaToken {
		return nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, err
	}

	arg2, err := p.expression(1)
	if err != nil {
		return nil, nil, err
	}

	if p.curr.Type == lexer.CommaToken {
		return nil, nil, &InvalidFunctionCallError{name}
	}

	if p.curr.Type != lexer.CloseParenToken {
		return nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, err
	}

	return arg1, arg2, nil
}

func (p *parser) function2ExpArg(name string) (Node, Node, error) {
	if p.curr.Type == lexer.CloseParenToken {
		return nil, nil, &InvalidFunctionCallError{name}
	}

	arg1, err := p.expression(1)
	if err != nil {
		return nil, nil, err
	}

	if p.curr.Type == lexer.CloseParenToken {
		return nil, nil, &InvalidFunctionCallError{name}
	}

	if p.curr.Type != lexer.CommaToken {
		return nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if p.next.Type != lexer.ExpressionToken {
		return nil, nil, &InvalidFunctionArgumentError{name, "expression"}
	}

	if err := p.advance2(); err != nil {
		return nil, nil, err
	}

	arg2, err := p.expression(1)
	if err != nil {
		return nil, nil, err
	}

	if p.curr.Type == lexer.CommaToken {
		return nil, nil, &InvalidFunctionCallError{}
	}

	if p.curr.Type != lexer.CloseParenToken {
		return nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, err
	}

	return arg1, arg2, nil
}

func (p *parser) function2MapArg(name string) (Node, Node, error) {
	if p.curr.Type == lexer.CloseParenToken {
		return nil, nil, &InvalidFunctionCallError{name}
	}

	if p.curr.Type != lexer.ExpressionToken {
		return nil, nil, &InvalidFunctionArgumentError{name, "expression"}
	}

	if err := p.advance(); err != nil {
		return nil, nil, err
	}

	arg1, err := p.expression(1)
	if err != nil {
		return nil, nil, err
	}

	if p.curr.Type == lexer.CloseParenToken {
		return nil, nil, &InvalidFunctionCallError{name}
	}

	if p.curr.Type != lexer.CommaToken {
		return nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, err
	}

	arg2, err := p.expression(1)
	if err != nil {
		return nil, nil, err
	}

	if p.curr.Type == lexer.CommaToken {
		return nil, nil, &InvalidFunctionCallError{}
	}

	if p.curr.Type != lexer.CloseParenToken {
		return nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, err
	}

	return arg1, arg2, nil
}

func (p *parser) function2To3Arg(name string) (Node, Node, Node, error) {
	if p.curr.Type == lexer.CloseParenToken {
		return nil, nil, nil, &InvalidFunctionCallError{name}
	}

	arg1, err := p.expression(1)
	if err != nil {
		return nil, nil, nil, err
	}

	if p.curr.Type == lexer.CloseParenToken {
		return nil, nil, nil, &InvalidFunctionCallError{name}
	}

	if p.curr.Type != lexer.CommaToken {
		return nil, nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, nil, err
	}

	arg2, err := p.expression(1)
	if err != nil {
		return nil, nil, nil, err
	}

	if p.curr.Type == lexer.CloseParenToken {
		if err := p.advance(); err != nil {
			return nil, nil, nil, err
		}

		return arg1, arg2, nil, nil
	}

	if p.curr.Type != lexer.CommaToken {
		return nil, nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, nil, err
	}

	arg3, err := p.expression(1)
	if err != nil {
		return nil, nil, nil, err
	}

	if p.curr.Type == lexer.CommaToken {
		return nil, nil, nil, &InvalidFunctionCallError{name}
	}

	if p.curr.Type != lexer.CloseParenToken {
		return nil, nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, nil, err
	}

	return arg1, arg2, arg3, nil
}

func (p *parser) function2To4Arg(name string) (Node, Node, Node, Node, error) {
	if p.curr.Type == lexer.CloseParenToken {
		return nil, nil, nil, nil, &InvalidFunctionCallError{name}
	}

	arg1, err := p.expression(1)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if p.curr.Type == lexer.CloseParenToken {
		return nil, nil, nil, nil, &InvalidFunctionCallError{name}
	}

	if p.curr.Type != lexer.CommaToken {
		return nil, nil, nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, nil, nil, err
	}

	arg2, err := p.expression(1)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if p.curr.Type == lexer.CloseParenToken {
		if err := p.advance(); err != nil {
			return nil, nil, nil, nil, err
		}

		return arg1, arg2, nil, nil, nil
	}

	if p.curr.Type != lexer.CommaToken {
		return nil, nil, nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, nil, nil, err
	}

	arg3, err := p.expression(1)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if p.curr.Type == lexer.CloseParenToken {
		if err := p.advance(); err != nil {
			return nil, nil, nil, nil, err
		}

		return arg1, arg2, arg3, nil, nil
	}

	if p.curr.Type != lexer.CommaToken {
		return nil, nil, nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, nil, nil, err
	}

	arg4, err := p.expression(1)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if p.curr.Type == lexer.CommaToken {
		return nil, nil, nil, nil, &InvalidFunctionCallError{name}
	}

	if p.curr.Type != lexer.CloseParenToken {
		return nil, nil, nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, nil, nil, err
	}

	return arg1, arg2, arg3, arg4, nil
}

func (p *parser) function3To4Arg(name string) (Node, Node, Node, Node, error) {
	if p.curr.Type == lexer.CloseParenToken {
		return nil, nil, nil, nil, &InvalidFunctionCallError{name}
	}

	arg1, err := p.expression(1)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if p.curr.Type == lexer.CloseParenToken {
		return nil, nil, nil, nil, &InvalidFunctionCallError{name}
	}

	if p.curr.Type != lexer.CommaToken {
		return nil, nil, nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, nil, nil, err
	}

	arg2, err := p.expression(1)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if p.curr.Type == lexer.CloseParenToken {
		return nil, nil, nil, nil, &InvalidFunctionCallError{name}
	}

	if p.curr.Type != lexer.CommaToken {
		return nil, nil, nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, nil, nil, err
	}

	arg3, err := p.expression(1)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if p.curr.Type == lexer.CloseParenToken {
		if err := p.advance(); err != nil {
			return nil, nil, nil, nil, err
		}

		return arg1, arg2, arg3, nil, nil
	}

	if p.curr.Type != lexer.CommaToken {
		return nil, nil, nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, nil, nil, err
	}

	arg4, err := p.expression(1)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if p.curr.Type == lexer.CommaToken {
		return nil, nil, nil, nil, &InvalidFunctionCallError{name}
	}

	if p.curr.Type != lexer.CloseParenToken {
		return nil, nil, nil, nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, nil, nil, nil, err
	}

	return arg1, arg2, arg3, arg4, nil
}

func (p *parser) functionNotNull() (Node, error) {
	if p.curr.Type == lexer.CloseParenToken {
		return nil, &InvalidFunctionCallError{"not_null"}
	}

	node, err := p.expression(1)
	if err != nil {
		return nil, err
	}

	if p.curr.Type == lexer.CloseParenToken {
		if err := p.advance(); err != nil {
			return nil, err
		}

		return &NotNullValueNode{
			Argument: node,
			Value:    nil,
		}, nil
	}

	if p.curr.Type != lexer.CommaToken {
		return nil, &unexpectedTokenError{p.curr.Value}
	}

	if err := p.advance(); err != nil {
		return nil, err
	}

	if p.next.Type == lexer.CloseParenToken {
		switch p.curr.Type {
		case lexer.JSONLiteralToken:
			value, err := parseJSONLiteral(p.curr.Value)
			if err != nil {
				return nil, err
			}

			if err := p.advance2(); err != nil {
				return nil, err
			}

			return &NotNullValueNode{
				Argument: node,
				Value:    value,
			}, nil
		case lexer.StringLiteralToken:
			value, err := parseStringLiteral(p.curr.Value)
			if err != nil {
				return nil, err
			}

			if err := p.advance2(); err != nil {
				return nil, err
			}

			return &NotNullValueNode{
				Argument: node,
				Value:    value,
			}, nil
		}
	}

	nodes := []Node{node}
	for {
		node, err := p.expression(1)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)

		if p.curr.Type == lexer.CommaToken {
			if err := p.advance(); err != nil {
				return nil, err
			}

			continue
		}

		if p.curr.Type == lexer.CloseParenToken {
			if err := p.advance(); err != nil {
				return nil, err
			}

			return &NotNullNode{
				Arguments: nodes,
			}, nil
		}

		return nil, &unexpectedTokenError{p.curr.Value}
	}
}

func (p *parser) functionVarArg(name string) ([]Node, error) {
	if p.curr.Type == lexer.CloseParenToken {
		return nil, &InvalidFunctionCallError{name}
	}

	var nodes []Node
	for {
		node, err := p.expression(1)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)

		if p.curr.Type == lexer.CommaToken {
			if err := p.advance(); err != nil {
				return nil, err
			}

			continue
		}

		if p.curr.Type == lexer.CloseParenToken {
			if err := p.advance(); err != nil {
				return nil, err
			}

			return nodes, nil
		}

		return nil, &unexpectedTokenError{p.curr.Value}
	}
}

func (p *parser) index(child Node) (Node, bool, error) {
	var haveStart bool
	start := 0
	if p.curr.Type == lexer.IntegerLiteralToken {
		var err error
		start, err = strconv.Atoi(p.curr.Value)
		if err != nil {
			return nil, false, &invalidIndexError{p.curr.Value}
		}

		if p.next.Type == lexer.CloseSqBraceToken {
			if err := p.advance2(); err != nil {
				return nil, false, err
			}

			if child == nil {
				if start >= 0 && start <= math.MaxUint8 {
					return SmallIndexCurrentNode{
						Value: uint8(start),
					}, false, nil
				}

				return &IndexCurrentNode{
					Value: start,
				}, false, nil
			}

			return &IndexNode{
				Child: child,
				Value: start,
			}, false, nil
		} else if p.next.Type == lexer.ColonToken {
			if err := p.advance2(); err != nil {
				return nil, false, err
			}
		} else {
			return nil, false, &unexpectedTokenError{p.next.Value}
		}

		haveStart = true
	} else if p.curr.Type == lexer.ColonToken {
		if err := p.advance(); err != nil {
			return nil, false, err
		}
	} else {
		return nil, false, &unexpectedTokenError{p.curr.Value}
	}

	var haveStop bool
	stop := math.MaxInt
	if p.curr.Type == lexer.IntegerLiteralToken {
		var err error
		stop, err = strconv.Atoi(p.curr.Value)
		if err != nil {
			return nil, false, &invalidIndexError{p.curr.Value}
		}

		if p.next.Type == lexer.CloseSqBraceToken {
			if err := p.advance2(); err != nil {
				return nil, false, err
			}

			if child == nil {
				return &SliceCurrentNode{
					Start: start,
					Stop:  stop,
				}, true, nil
			}

			return &SliceNode{
				Child: child,
				Start: start,
				Stop:  stop,
			}, true, nil
		} else if p.next.Type == lexer.ColonToken {
			if err := p.advance2(); err != nil {
				return nil, false, err
			}
		} else {
			return nil, false, &unexpectedTokenError{p.next.Value}
		}

		haveStop = true
	} else if p.curr.Type == lexer.CloseSqBraceToken {
		if err := p.advance(); err != nil {
			return nil, false, err
		}

		if child == nil {
			return &SliceCurrentNode{
				Start: start,
				Stop:  math.MaxInt,
			}, true, nil
		}

		return &SliceNode{
			Child: child,
			Start: start,
			Stop:  math.MaxInt,
		}, true, nil
	} else if p.curr.Type == lexer.ColonToken {
		if err := p.advance(); err != nil {
			return nil, false, err
		}
	} else {
		return nil, false, &unexpectedTokenError{p.curr.Value}
	}

	step := 1
	if p.curr.Type == lexer.IntegerLiteralToken {
		if p.next.Type != lexer.CloseSqBraceToken {
			return nil, false, &unexpectedTokenError{p.next.Value}
		}

		var err error
		step, err = strconv.Atoi(p.curr.Value)
		if err != nil {
			return nil, false, &invalidIndexError{p.curr.Value}
		}

		if step == 0 {
			return nil, false, &InvalidSliceStepError{}
		}

		if step < 0 {
			if !haveStart {
				start = math.MaxInt
			}

			if !haveStop {
				stop = math.MinInt
			}
		}

		if err := p.advance2(); err != nil {
			return nil, false, err
		}
	} else if p.curr.Type == lexer.CloseSqBraceToken {
		if err := p.advance(); err != nil {
			return nil, false, err
		}
	}

	if child == nil {
		if step == 1 {
			return &SliceCurrentNode{
				Start: start,
				Stop:  stop,
			}, true, nil
		}

		return &SliceStepCurrentNode{
			Start: start,
			Stop:  stop,
			Step:  step,
		}, true, nil
	}

	if step == 1 {
		return &SliceNode{
			Child: child,
			Start: start,
			Stop:  stop,
		}, true, nil
	}

	return &SliceStepNode{
		Child: child,
		Start: start,
		Stop:  stop,
		Step:  step,
	}, true, nil
}

func (p *parser) let() (Node, error) {
	variables := make(map[string]Node)
	for {
		if p.curr.Type != lexer.VariableToken {
			return nil, &unexpectedTokenError{p.curr.Value}
		}

		if p.next.Type != lexer.AssignToken {
			return nil, &unexpectedTokenError{p.next.Value}
		}

		variable := p.curr.Value

		if err := p.advance2(); err != nil {
			return nil, err
		}

		node, err := p.expression(1)
		if err != nil {
			return nil, err
		}

		variables[variable] = node

		if p.curr.Type == lexer.InToken {
			if err := p.advance(); err != nil {
				return nil, err
			}

			break
		}

		if p.curr.Type != lexer.CommaToken {
			return nil, &unexpectedTokenError{p.curr.Value}
		}

		if err := p.advance(); err != nil {
			return nil, err
		}
	}

	child, err := p.expression(1)
	if err != nil {
		return nil, err
	}

	return &DefineVariables{
		Variables: variables,
		Child:     child,
	}, nil
}

func (p *parser) parse() (Node, error) {
	node, err := p.expression(1)
	if err != nil {
		return nil, err
	}

	if p.curr.Type != lexer.EndToken {
		return nil, &unexpectedTokenError{p.curr.Value}
	}

	return node, nil
}

func (p *parser) primaryExpression() (Node, error) {
	var node Node
	var err error
	switch p.curr.Type {
	case lexer.AddToken:
		if err := p.advance(); err != nil {
			return nil, err
		}

		child, err := p.expression(precedence(lexer.AddToken))
		if err != nil {
			return nil, err
		}

		node = &AssertNumberNode{
			Child: child,
		}
	case lexer.ArrayWildcardToken:
		if err := p.advance(); err != nil {
			return nil, err
		}

		child, err := p.projection(precedence(lexer.ObjectWildcardToken))
		if err != nil {
			return nil, err
		}

		if child == nil {
			node = PruneArrayCurrentNode{}
		} else {
			node = &ProjectArrayCurrentNode{
				Child: child,
			}
		}
	case lexer.AsteriskToken:
		if err := p.advance(); err != nil {
			return nil, err
		}

		child, err := p.projection(precedence(lexer.ObjectWildcardToken))
		if err != nil {
			return nil, err
		}

		if child == nil {
			node = ObjectValuesCurrentNode{}
		} else {
			node = &ProjectObjectCurrentNode{
				Child: child,
			}
		}
	case lexer.CurrentToken:
		node = CurrentNode{}

		if err := p.advance(); err != nil {
			return nil, err
		}
	case lexer.FilterToken:
		if err := p.advance(); err != nil {
			return nil, err
		}

		filter, err := p.filter()
		if err != nil {
			return nil, err
		}

		child, err := p.projection(precedence(lexer.FilterToken))
		if err != nil {
			return nil, err
		}

		if child == nil {
			node = &FilterCurrentNode{
				Filter: filter,
			}
		} else {
			node = &FilterAndProjectCurrentNode{
				Filter: filter,
				Child:  child,
			}
		}
	case lexer.FlattenToken:
		if err := p.advance(); err != nil {
			return nil, err
		}

		child, err := p.projection(precedence(lexer.FlattenToken))
		if err != nil {
			return nil, err
		}

		if child == nil {
			node = FlattenCurrentNode{}
		} else {
			node = &FlattenAndProjectCurrentNode{
				Child: child,
			}
		}
	case lexer.JSONLiteralToken:
		value, err := parseJSONLiteral(p.curr.Value)
		if err != nil {
			return nil, err
		}

		switch value := value.(type) {
		case bool:
			node = BoolNode{
				Value: value,
			}
		case nil:
			node = NullNode{}
		default:
			node = &ValueNode{
				Value: value,
			}
		}

		if err := p.advance(); err != nil {
			return nil, err
		}
	case lexer.LetToken:
		if err := p.advance(); err != nil {
			return nil, err
		}

		node, err = p.let()
		if err != nil {
			return nil, err
		}
	case lexer.NotToken:
		if err := p.advance(); err != nil {
			return nil, err
		}

		node, err = p.expression(precedence(lexer.NotToken))
		if err != nil {
			return nil, err
		}

		node = &NotNode{
			Child: node,
		}
	case lexer.OpenParenToken:
		if err := p.advance(); err != nil {
			return nil, err
		}

		node, err = p.expression(1)
		if err != nil {
			return nil, err
		}

		if p.curr.Type != lexer.CloseParenToken {
			return nil, &unexpectedTokenError{p.curr.Value}
		}

		if err := p.advance(); err != nil {
			return nil, err
		}
	case lexer.OpenBraceToken:
		if err := p.advance(); err != nil {
			return nil, err
		}

		node, err = p.selectObject(nil)
		if err != nil {
			return nil, err
		}
	case lexer.OpenSqBraceToken:
		if err := p.advance(); err != nil {
			return nil, err
		}

		if p.curr.Type == lexer.IntegerLiteralToken || p.curr.Type == lexer.ColonToken {
			var project bool
			node, project, err = p.index(nil)
			if err != nil {
				return nil, err
			}

			if project {
				right, err := p.projection(precedence(lexer.OpenSqBraceToken))
				if err != nil {
					return nil, err
				}

				if right != nil {
					node = &ProjectArrayNode{
						Left:  node,
						Right: right,
					}
				}
			}
		} else {
			node, err = p.selectArray(nil)
			if err != nil {
				return nil, err
			}
		}
	case lexer.QuotedIdentifierToken:
		value, err := parseQuotedIdentifier(p.curr.Value)
		if err != nil {
			return nil, err
		}

		node = &FieldNode{
			Value: value,
		}

		if err := p.advance(); err != nil {
			return nil, err
		}
	case lexer.RootToken:
		node = RootNode{}

		if err := p.advance(); err != nil {
			return nil, err
		}
	case lexer.StringLiteralToken:
		value, err := parseStringLiteral(p.curr.Value)
		if err != nil {
			return nil, err
		}

		node = &ValueNode{
			Value: value,
		}

		if err := p.advance(); err != nil {
			return nil, err
		}
	case lexer.SubtractToken:
		if err := p.advance(); err != nil {
			return nil, err
		}

		child, err := p.expression(precedence(lexer.SubtractToken))
		if err != nil {
			return nil, err
		}

		return &NegateNode{
			Child: child,
		}, nil
	case lexer.UnquotedIdentifierToken:
		if p.next.Type == lexer.OpenParenToken {
			node, err = p.function()
			if err != nil {
				return nil, err
			}
		} else {
			node = &FieldNode{
				Value: p.curr.Value,
			}

			if err := p.advance(); err != nil {
				return nil, err
			}
		}
	case lexer.VariableToken:
		node = &VariableNode{
			Name: p.curr.Value,
		}

		if err := p.advance(); err != nil {
			return nil, err
		}
	default:
		return nil, &unexpectedTokenError{p.curr.Value}
	}

	return node, nil
}

func (p *parser) projection(prec int) (Node, error) {
	var node Node
	var err error
	switch p.curr.Type {
	case lexer.DotToken:
		switch p.next.Type {
		case lexer.ArrayWildcardToken:
			if err := p.advance2(); err != nil {
				return nil, err
			}

			node = &SelectArraySingleCurrentNode{
				Field: ObjectValuesCurrentNode{},
			}
		case lexer.OpenBraceToken:
			if err := p.advance2(); err != nil {
				return nil, err
			}

			node, err = p.selectObject(nil)
			if err != nil {
				return nil, err
			}
		case lexer.OpenSqBraceToken:
			if err := p.advance2(); err != nil {
				return nil, err
			}

			node, err = p.selectArray(nil)
			if err != nil {
				return nil, err
			}
		case lexer.QuotedIdentifierToken,
			lexer.UnquotedIdentifierToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			node, err = p.expression(prec)
			if err != nil {
				return nil, err
			}
		default:
			return nil, &unexpectedTokenError{p.curr.Value}
		}
	case lexer.FilterToken:
		if err := p.advance(); err != nil {
			return nil, err
		}

		filter, err := p.filter()
		if err != nil {
			return nil, err
		}

		node = &FilterCurrentNode{
			Filter: filter,
		}
	case lexer.ObjectWildcardToken:
		if p.next.Type == lexer.EndToken {
			if err := p.advance(); err != nil {
				return nil, err
			}

			node = ObjectValuesCurrentNode{}
		} else {
			p.setCurrent(lexer.Token{
				Type:  lexer.AsteriskToken,
				Value: p.curr.Value[1:],
			})

			node, err = p.expression(prec)
			if err != nil {
				return nil, err
			}
		}
	case lexer.OpenSqBraceToken:
		if err := p.advance(); err != nil {
			return nil, err
		}

		node, _, err = p.index(nil)
		if err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}

	newPrec := precedence(p.curr.Type)
	for newPrec > prec {
		switch p.curr.Type {
		case lexer.DotToken:
			switch p.next.Type {
			case lexer.ArrayWildcardToken:
				if err := p.advance2(); err != nil {
					return nil, err
				}

				node = &SelectArraySingleNode{
					Child: node,
					Field: ObjectValuesCurrentNode{},
				}
			case lexer.OpenBraceToken:
				if err := p.advance2(); err != nil {
					return nil, err
				}

				node, err = p.selectObject(node)
				if err != nil {
					return nil, err
				}
			case lexer.OpenSqBraceToken:
				if err := p.advance2(); err != nil {
					return nil, err
				}

				node, err = p.selectArray(node)
				if err != nil {
					return nil, err
				}
			case lexer.QuotedIdentifierToken,
				lexer.UnquotedIdentifierToken:
				if err := p.advance(); err != nil {
					return nil, err
				}

				node, err = p.expression(newPrec)
				if err != nil {
					return nil, err
				}
			default:
				return nil, &unexpectedTokenError{p.curr.Value}
			}
		case lexer.FilterToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			filter, err := p.filter()
			if err != nil {
				return nil, err
			}

			node = &FilterNode{
				Child:  node,
				Filter: filter,
			}
		case lexer.ObjectWildcardToken:
			if p.curr.Type == lexer.EndToken {
				if err := p.advance(); err != nil {
					return nil, err
				}

				node = &ObjectValuesNode{
					Child: node,
				}
			} else {
				p.setCurrent(lexer.Token{
					Type:  lexer.AsteriskToken,
					Value: p.curr.Value[1:],
				})

				right, err := p.expression(newPrec)
				if err != nil {
					return nil, err
				}

				node = &ProjectObjectNode{
					Left:  node,
					Right: right,
				}
			}
		case lexer.OpenSqBraceToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			node, _, err = p.index(node)
			if err != nil {
				return nil, err
			}
		default:
			return nil, &unexpectedTokenError{p.curr.Value}
		}

		newPrec = precedence(p.curr.Type)
	}

	return node, nil
}

func (p *parser) selectArray(child Node) (Node, error) {
	var fields []Node
	for {
		field, err := p.expression(1)
		if err != nil {
			return nil, err
		}

		switch p.curr.Type {
		case lexer.CommaToken:
			fields = append(fields, field)

			if err := p.advance(); err != nil {
				return nil, err
			}
		case lexer.CloseSqBraceToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			if len(fields) == 0 {
				if child == nil {
					return &SelectArraySingleCurrentNode{
						Field: field,
					}, nil
				}

				return &SelectArraySingleNode{
					Child: child,
					Field: field,
				}, nil
			}

			fields = append(fields, field)

			if child == nil {
				return &SelectArrayCurrentNode{
					Fields: fields,
				}, nil
			}

			return &SelectArrayNode{
				Child:  child,
				Fields: fields,
			}, nil
		default:
			return nil, &unexpectedTokenError{p.curr.Value}
		}
	}
}

func (p *parser) selectObject(child Node) (Node, error) {
	fields := make(map[string]Node)
	for {
		var key string
		switch p.curr.Type {
		case lexer.QuotedIdentifierToken:
			var err error
			key, err = parseQuotedIdentifier(p.curr.Value)
			if err != nil {
				return nil, err
			}
		case lexer.UnquotedIdentifierToken:
			key = p.curr.Value
		}

		if p.next.Type != lexer.ColonToken {
			return nil, &unexpectedTokenError{p.next.Value}
		}

		if err := p.advance2(); err != nil {
			return nil, err
		}

		field, err := p.expression(1)
		if err != nil {
			return nil, err
		}

		switch p.curr.Type {
		case lexer.CommaToken:
			fields[key] = field

			if err := p.advance(); err != nil {
				return nil, err
			}
		case lexer.CloseBraceToken:
			if err := p.advance(); err != nil {
				return nil, err
			}

			if len(fields) == 0 {
				if child == nil {
					return &SelectObjectSingleCurrentNode{
						Key:   key,
						Field: field,
					}, nil
				}

				return &SelectObjectSingleNode{
					Child: child,
					Key:   key,
					Field: field,
				}, nil
			}

			fields[key] = field

			if child == nil {
				return &SelectObjectCurrentNode{
					Fields: fields,
				}, nil
			}

			return &SelectObjectNode{
				Child:  child,
				Fields: fields,
			}, nil
		}
	}
}

func (p *parser) setCurrent(tok lexer.Token) {
	p.curr = tok
}

func parseJSONLiteral(s string) (any, error) {
	v := strings.ReplaceAll(s[1:len(s)-1], "\\`", "`")
	if len(v) == 0 {
		return nil, &invalidJSONLiteralError{s}
	}

	switch v[0] {
	case '"':
		var s string
		if err := json.Unmarshal([]byte(v), &s); err == nil {
			return s, nil
		}
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		var n json.Number
		if err := json.Unmarshal([]byte(v), &n); err == nil {
			return n, nil
		}
	case 'f':
		if v == "false" {
			return false, nil
		}
	case 'n':
		if v == "null" {
			return nil, nil
		}
	case 't':
		if v == "true" {
			return true, nil
		}
	}

	r := strings.NewReader(v)
	d := json.NewDecoder(r)
	d.UseNumber()

	var a any
	if err := d.Decode(&a); err != nil {
		return nil, &invalidJSONLiteralError{s}
	}

	if _, err := d.Token(); err == nil || !errors.Is(err, io.EOF) {
		return nil, &invalidJSONLiteralError{s}
	}

	return a, nil
}

func parseQuotedIdentifier(s string) (string, error) {
	v := s[1 : len(s)-1]
	i := strings.IndexByte(v, '\\')
	if i == -1 || i+1 == len(v) {
		return v, nil
	}

	var b strings.Builder
	b.Grow(len(v))
	b.WriteString(v[:i])

	v = v[i+1:]
	for {
		switch v[0] {
		case '"':
			b.WriteByte('"')
			v = v[1:]
		case '/':
			b.WriteByte('/')
			v = v[1:]
		case '\\':
			b.WriteByte('\\')
			v = v[1:]
		case 'b':
			b.WriteByte('\b')
			v = v[1:]
		case 'f':
			b.WriteByte('\f')
			v = v[1:]
		case 'n':
			b.WriteByte('\n')
			v = v[1:]
		case 'r':
			b.WriteByte('\r')
			v = v[1:]
		case 't':
			b.WriteByte('\t')
			v = v[1:]
		case 'u':
			if len(v) < 5 {
				return "", &invalidQuotedStringError{s}
			}

			var r rune
			for _, c := range v[1:5] {
				if c >= '0' && c <= '9' {
					r = r*16 + rune(c-'0')
				} else if c >= 'a' && c <= 'f' {
					r = r*16 + rune(c-'a'+10)
				} else if c >= 'A' && c <= 'F' {
					r = r*16 + rune(c-'A'+10)
				} else {
					return "", &invalidQuotedStringError{s}
				}
			}

			v = v[5:]

			if utf16.IsSurrogate(r) {
				if len(v) < 6 {
					return "", &invalidQuotedStringError{s}
				}

				if v[0] != '\\' && v[1] != 'u' {
					return "", &invalidQuotedStringError{s}
				}

				var r2 rune
				for _, c := range v[2:6] {
					if c >= 0 && c <= '9' {
						r2 = r2*16 + rune(c-'0')
					} else if c >= 'a' && c <= 'f' {
						r2 = r2*16 + rune(c-'a'+10)
					} else if c >= 'A' && c <= 'F' {
						r2 = r2*16 + rune(c-'A'+10)
					} else {
						return "", &invalidQuotedStringError{s}
					}
				}

				r = utf16.DecodeRune(r, r2)
				v = v[6:]
			}

			b.WriteRune(r)
		default:
			return "", &invalidQuotedStringError{s}
		}

		i := strings.IndexByte(v, '\\')
		if i == -1 || i+1 == len(v) {
			b.WriteString(v)

			return b.String(), nil
		}

		b.WriteString(v[:i])
		v = v[i+1:]
	}
}

func parseStringLiteral(s string) (string, error) {
	v := s[1 : len(s)-1]
	i := strings.IndexByte(v, '\\')
	if i == -1 || i+1 == len(v) {
		return v, nil
	}

	var b strings.Builder
	b.Grow(len(v))
	b.WriteString(v[:i])

	v = v[i+1:]
	for {
		switch v[0] {
		case '\'':
			b.WriteByte('\'')
		case '\\':
			b.WriteByte('\\')
		default:
			b.WriteByte('\\')
			b.WriteByte(v[0])
		}

		v = v[1:]
		i := strings.IndexByte(v, '\\')
		if i == -1 || i+1 == len(v) {
			b.WriteString(v)

			return b.String(), nil
		}

		b.WriteString(v[:i])
		v = v[i+1:]
	}
}
