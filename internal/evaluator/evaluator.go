package evaluator

import (
	"math"
	"reflect"

	"github.com/woodsbury/jmespath/internal/parser"
)

func Evaluate(node parser.Node, data any) (any, error) {
	e := evaluator{
		root: data,
	}

	return e.evaluate(node, data, nil)
}

type evaluator struct {
	root any
}

func (e *evaluator) evaluate(node parser.Node, current any, variables *variableScope) (any, error) {
	switch node := node.(type) {
	case *parser.AbsNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return abs(arg)
	case *parser.AddNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		right, err := e.evaluate(node.Right, current, variables)
		if err != nil {
			return nil, err
		}

		return add(left, right)
	case *parser.AndNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		if !isTrue(left) {
			return left, nil
		}

		return e.evaluate(node.Right, current, variables)
	case *parser.AssertNumberNode:
		child, err := e.evaluate(node.Child, current, variables)
		if err != nil {
			return nil, err
		}

		if isNumber(child) {
			return child, nil
		}

		return nil, nil
	case *parser.AvgNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return avg(arg)
	case parser.BoolNode:
		return node.Value, nil
	case *parser.CeilNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return ceil(arg)
	case *parser.ContainsNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		return contains(arg1, arg2)
	case parser.CurrentNode:
		return current, nil
	case *parser.DefineVariables:
		results := make(map[string]any, len(node.Variables))
		for name, node := range node.Variables {
			result, err := e.evaluate(node, current, variables)
			if err != nil {
				return nil, err
			}

			results[name] = result
		}

		return e.evaluate(node.Child, current, variables.new(results))
	case *parser.DivideNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		right, err := e.evaluate(node.Right, current, variables)
		if err != nil {
			return nil, err
		}

		return divide(left, right)
	case *parser.EndsWithNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		return endsWith(arg1, arg2)
	case *parser.EqualNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		right, err := e.evaluate(node.Right, current, variables)
		if err != nil {
			return nil, err
		}

		return equal(left, right), nil
	case *parser.FieldNode:
		return field(node.Value, current), nil
	case *parser.FilterNode:
		child, err := e.evaluate(node.Child, current, variables)
		if err != nil {
			return nil, err
		}

		return e.filter(child, node.Filter, variables)
	case *parser.FilterAndProjectNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		return e.filterAndProjectArray(left, node.Filter, node.Right, variables)
	case *parser.FilterAndProjectCurrentNode:
		return e.filterAndProjectArray(current, node.Filter, node.Child, variables)
	case *parser.FilterCurrentNode:
		return e.filter(current, node.Filter, variables)
	case *parser.FindFirstNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		return findFirst(arg1, arg2)
	case *parser.FindFirstBetweenNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		arg3, err := e.evaluate(node.Arguments[2], current, variables)
		if err != nil {
			return nil, err
		}

		arg4, err := e.evaluate(node.Arguments[3], current, variables)
		if err != nil {
			return nil, err
		}

		return findFirstBetween(arg1, arg2, arg3, arg4)
	case *parser.FindFirstFromNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		arg3, err := e.evaluate(node.Arguments[2], current, variables)
		if err != nil {
			return nil, err
		}

		return findFirstFrom(arg1, arg2, arg3)
	case *parser.FindLastNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		return findLast(arg1, arg2)
	case *parser.FindLastBetweenNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		arg3, err := e.evaluate(node.Arguments[2], current, variables)
		if err != nil {
			return nil, err
		}

		arg4, err := e.evaluate(node.Arguments[3], current, variables)
		if err != nil {
			return nil, err
		}

		return findLastBetween(arg1, arg2, arg3, arg4)
	case *parser.FindLastFromNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		arg3, err := e.evaluate(node.Arguments[2], current, variables)
		if err != nil {
			return nil, err
		}

		return findLastFrom(arg1, arg2, arg3)
	case *parser.FlattenNode:
		child, err := e.evaluate(node.Child, current, variables)
		if err != nil {
			return nil, err
		}

		return flatten(child), nil
	case *parser.FlattenAndProjectNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		return e.flattenAndProjectArray(left, node.Right, variables)
	case *parser.FlattenAndProjectCurrentNode:
		return e.flattenAndProjectArray(current, node.Child, variables)
	case parser.FlattenCurrentNode:
		return flatten(current), nil
	case *parser.FloorNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return floor(arg)
	case *parser.FromItemsNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return fromItems(arg)
	case *parser.GreaterNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		right, err := e.evaluate(node.Right, current, variables)
		if err != nil {
			return nil, err
		}

		return greater(left, right), nil
	case *parser.GreaterOrEqualNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		right, err := e.evaluate(node.Right, current, variables)
		if err != nil {
			return nil, err
		}

		return greaterOrEqual(left, right), nil
	case *parser.GroupByNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		return e.groupBy(arg1, node.Arguments[1], variables)
	case *parser.IndexNode:
		child, err := e.evaluate(node.Child, current, variables)
		if err != nil {
			return nil, err
		}

		return index(child, node.Value), nil
	case *parser.IndexCurrentNode:
		return index(current, node.Value), nil
	case *parser.IntegerDivideNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		right, err := e.evaluate(node.Right, current, variables)
		if err != nil {
			return nil, err
		}

		return integerDivide(left, right)
	case *parser.ItemsNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return items(arg)
	case *parser.JoinNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		return join(arg1, arg2)
	case *parser.KeysNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return keys(arg)
	case *parser.LengthNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return length(arg)
	case *parser.LessNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		right, err := e.evaluate(node.Right, current, variables)
		if err != nil {
			return nil, err
		}

		return less(left, right), nil
	case *parser.LessOrEqualNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		right, err := e.evaluate(node.Right, current, variables)
		if err != nil {
			return nil, err
		}

		return lessOrEqual(left, right), nil
	case *parser.LowerNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return lower(arg)
	case *parser.MapNode:
		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		return e.mapArray(arg2, node.Arguments[0], variables)
	case *parser.MaxNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return arrayMax(arg)
	case *parser.MaxByNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		return e.arrayMaxBy(arg1, node.Arguments[1], variables)
	case *parser.MergeNode:
		result := make(map[string]any)
		for _, arg := range node.Arguments {
			value, err := e.evaluate(arg, current, variables)
			if err != nil {
				return nil, err
			}

			m, ok := value.(map[string]any)
			if !ok {
				return nil, &InvalidTypeError{
					got:  reflect.TypeOf(value),
					want: "object",
				}
			}

			for k, v := range m {
				result[k] = v
			}
		}

		return result, nil
	case *parser.MinNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return arrayMin(arg)
	case *parser.MinByNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		return e.arrayMinBy(arg1, node.Arguments[1], variables)
	case *parser.ModuloNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		right, err := e.evaluate(node.Right, current, variables)
		if err != nil {
			return nil, err
		}

		return modulo(left, right)
	case *parser.MultiplyNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		right, err := e.evaluate(node.Right, current, variables)
		if err != nil {
			return nil, err
		}

		return multiply(left, right)
	case *parser.NegateNode:
		child, err := e.evaluate(node.Child, current, variables)
		if err != nil {
			return nil, err
		}

		if f, ok := toFloat(child); ok {
			return -f, nil
		}

		d, ok := toDecimal(child)
		if !ok {
			return nil, nil
		}

		if d.IsZero() {
			return d, nil
		}

		return d.Neg(), nil
	case *parser.NotNode:
		child, err := e.evaluate(node.Child, current, variables)
		if err != nil {
			return nil, err
		}

		return !isTrue(child), nil
	case *parser.NotEqualNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		right, err := e.evaluate(node.Right, current, variables)
		if err != nil {
			return nil, err
		}

		return !equal(left, right), nil
	case *parser.NotNullNode:
		for _, arg := range node.Arguments {
			result, err := e.evaluate(arg, current, variables)
			if err != nil {
				return nil, err
			}

			if result != nil {
				return result, nil
			}
		}

		return nil, nil
	case parser.NullNode:
		return nil, nil
	case *parser.ObjectValuesNode:
		child, err := e.evaluate(node.Child, current, variables)
		if err != nil {
			return nil, err
		}

		return objectValues(child), nil
	case parser.ObjectValuesCurrentNode:
		return objectValues(current), nil
	case *parser.OrNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		if isTrue(left) {
			return left, nil
		}

		return e.evaluate(node.Right, current, variables)
	case *parser.PadLeftNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		arg3, err := e.evaluate(node.Arguments[2], current, variables)
		if err != nil {
			return nil, err
		}

		return padLeft(arg1, arg2, arg3)
	case *parser.PadRightNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		arg3, err := e.evaluate(node.Arguments[2], current, variables)
		if err != nil {
			return nil, err
		}

		return padRight(arg1, arg2, arg3)
	case *parser.PadSpaceLeftNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		return padSpaceLeft(arg1, arg2)
	case *parser.PadSpaceRightNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		return padSpaceRight(arg1, arg2)
	case *parser.PipeNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		return e.evaluate(node.Right, left, variables)
	case *parser.PipeFieldNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		return field(node.Right, left), nil
	case *parser.ProjectArrayNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		if _, ok := left.(string); ok && isSliceNode(node.Left) {
			return e.evaluate(node.Right, left, variables)
		}

		return e.projectArray(left, node.Right, variables)
	case *parser.ProjectArrayCurrentNode:
		return e.projectArray(current, node.Child, variables)
	case *parser.ProjectObjectNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		return e.projectObject(left, node.Right, variables)
	case *parser.ProjectObjectCurrentNode:
		return e.projectObject(current, node.Child, variables)
	case *parser.PruneArrayNode:
		child, err := e.evaluate(node.Child, current, variables)
		if err != nil {
			return nil, err
		}

		return pruneArray(child), nil
	case parser.PruneArrayCurrentNode:
		return pruneArray(current), nil
	case *parser.ReplaceNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		arg3, err := e.evaluate(node.Arguments[2], current, variables)
		if err != nil {
			return nil, err
		}

		return replace(arg1, arg2, arg3)
	case *parser.ReplaceCountNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		arg3, err := e.evaluate(node.Arguments[2], current, variables)
		if err != nil {
			return nil, err
		}

		arg4, err := e.evaluate(node.Arguments[3], current, variables)
		if err != nil {
			return nil, err
		}

		return replaceCount(arg1, arg2, arg3, arg4)
	case *parser.ReverseNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return reverse(arg)
	case parser.RootNode:
		return e.root, nil
	case *parser.SelectArrayNode:
		child, err := e.evaluate(node.Child, current, variables)
		if err != nil {
			return nil, err
		}

		if child == nil {
			return nil, nil
		}

		results := make([]any, len(node.Fields))
		for i, field := range node.Fields {
			result, err := e.evaluate(field, child, variables)
			if err != nil {
				return nil, err
			}

			results[i] = result
		}

		return results, nil
	case *parser.SelectArrayCurrentNode:
		if current == nil {
			return nil, nil
		}

		results := make([]any, len(node.Fields))
		for i, field := range node.Fields {
			result, err := e.evaluate(field, current, variables)
			if err != nil {
				return nil, err
			}

			results[i] = result
		}

		return results, nil
	case *parser.SelectArraySingleNode:
		child, err := e.evaluate(node.Child, current, variables)
		if err != nil {
			return nil, err
		}

		if child == nil {
			return nil, nil
		}

		result, err := e.evaluate(node.Field, child, variables)
		if err != nil {
			return nil, err
		}

		return []any{result}, nil
	case *parser.SelectArraySingleCurrentNode:
		result, err := e.evaluate(node.Field, current, variables)
		if err != nil {
			return nil, err
		}

		return []any{result}, nil
	case *parser.SelectObjectNode:
		child, err := e.evaluate(node.Child, current, variables)
		if err != nil {
			return nil, err
		}

		if child == nil {
			return nil, nil
		}

		results := make(map[string]any, len(node.Fields))
		for key, field := range node.Fields {
			result, err := e.evaluate(field, child, variables)
			if err != nil {
				return nil, err
			}

			results[key] = result
		}

		return results, nil
	case *parser.SelectObjectCurrentNode:
		if current == nil {
			return nil, nil
		}

		results := make(map[string]any, len(node.Fields))
		for key, field := range node.Fields {
			result, err := e.evaluate(field, current, variables)
			if err != nil {
				return nil, err
			}

			results[key] = result
		}

		return results, nil
	case *parser.SelectObjectSingleNode:
		child, err := e.evaluate(node.Child, current, variables)
		if err != nil {
			return nil, err
		}

		if child == nil {
			return nil, nil
		}

		result, err := e.evaluate(node.Field, child, variables)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			node.Key: result,
		}, nil
	case *parser.SelectObjectSingleCurrentNode:
		result, err := e.evaluate(node.Field, current, variables)
		if err != nil {
			return nil, err
		}

		return map[string]any{
			node.Key: result,
		}, nil
	case *parser.SliceNode:
		child, err := e.evaluate(node.Child, current, variables)
		if err != nil {
			return nil, err
		}

		return slice(child, node.Start, node.Stop), nil
	case *parser.SliceCurrentNode:
		return slice(current, node.Start, node.Stop), nil
	case *parser.SliceStepNode:
		child, err := e.evaluate(node.Child, current, variables)
		if err != nil {
			return nil, err
		}

		return sliceStep(child, node.Start, node.Stop, node.Step), nil
	case *parser.SliceStepCurrentNode:
		return sliceStep(current, node.Start, node.Stop, node.Step), nil
	case parser.SmallIndexCurrentNode:
		return index(current, int(node.Value)), nil
	case *parser.SortNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return sortArray(arg)
	case *parser.SortByNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		return e.sortArrayBy(arg1, node.Arguments[1], variables)
	case *parser.SplitNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		return split(arg1, arg2)
	case *parser.SplitCountNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		arg3, err := e.evaluate(node.Arguments[2], current, variables)
		if err != nil {
			return nil, err
		}

		return splitCount(arg1, arg2, arg3)
	case *parser.StartsWithNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		return startsWith(arg1, arg2)
	case *parser.SubtractNode:
		left, err := e.evaluate(node.Left, current, variables)
		if err != nil {
			return nil, err
		}

		right, err := e.evaluate(node.Right, current, variables)
		if err != nil {
			return nil, err
		}

		return subtract(left, right)
	case *parser.SumNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return sum(arg)
	case *parser.ToArrayNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return toArray(arg), nil
	case *parser.ToNumberNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return toNumber(arg), nil
	case *parser.ToStringNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return toString(arg)
	case *parser.TrimNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		return trim(arg1, arg2)
	case *parser.TrimLeftNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		return trimLeft(arg1, arg2)
	case *parser.TrimRightNode:
		arg1, err := e.evaluate(node.Arguments[0], current, variables)
		if err != nil {
			return nil, err
		}

		arg2, err := e.evaluate(node.Arguments[1], current, variables)
		if err != nil {
			return nil, err
		}

		return trimRight(arg1, arg2)
	case *parser.TrimSpaceNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return trimSpace(arg)
	case *parser.TrimSpaceLeftNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return trimSpaceLeft(arg)
	case *parser.TrimSpaceRightNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return trimSpaceRight(arg)
	case *parser.TypeNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return typeName(arg)
	case *parser.UpperNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return upper(arg)
	case *parser.ValueNode:
		return node.Value, nil
	case *parser.ValuesNode:
		arg, err := e.evaluate(node.Argument, current, variables)
		if err != nil {
			return nil, err
		}

		return values(arg)
	case *parser.VariableNode:
		value, ok := variables.get(node.Name)
		if !ok {
			return nil, &UndefinedVariableError{
				Variable: node.Name,
			}
		}

		return value, nil
	case *parser.ZipNode:
		count := math.MaxInt
		values := make([][]any, len(node.Arguments))
		for i, arg := range node.Arguments {
			value, err := e.evaluate(arg, current, variables)
			if err != nil {
				return nil, err
			}

			a, ok := value.([]any)
			if !ok {
				return nil, &InvalidTypeError{
					got:  reflect.TypeOf(value),
					want: "array",
				}
			}

			if l := len(a); l < count {
				count = l
			}

			values[i] = a
		}

		results := make([]any, count)
		for i := 0; i < count; i++ {
			result := make([]any, len(values))
			for j, value := range values {
				result[j] = value[i]
			}

			results[i] = result
		}

		return results, nil
	}

	return nil, &unexpectedOperationError{reflect.TypeOf(node)}
}
