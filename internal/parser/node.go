package parser

import (
	"encoding/json"
	"strconv"
)

type Node interface {
	ResultTypes() Types
	String() string
}

type AbsNode struct {
	Argument Node
}

func (n *AbsNode) ResultTypes() Types {
	return TypeNumber
}

func (n *AbsNode) String() string {
	return "Abs"
}

func (n *AbsNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type AddNode struct {
	Left  Node
	Right Node
}

func (n *AddNode) ResultTypes() Types {
	return TypeNumber
}

func (n *AddNode) String() string {
	return "Add"
}

func (n *AddNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type AndNode struct {
	Left  Node
	Right Node
}

func (n *AndNode) ResultTypes() Types {
	return TypeAny
}

func (n *AndNode) String() string {
	return "And"
}

func (n *AndNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type AssertNumberNode struct {
	Child Node
}

func (n *AssertNumberNode) ResultTypes() Types {
	return TypeNumber
}

func (n *AssertNumberNode) String() string {
	return "AssertNumber"
}

func (n *AssertNumberNode) Walk(v Visitor) {
	v.Visit(n.Child)
}

type AvgNode struct {
	Argument Node
}

func (n AvgNode) ResultTypes() Types {
	return TypeNumber
}

func (n *AvgNode) String() string {
	return "Avg"
}

func (n *AvgNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type BoolNode struct {
	Value bool
}

func (n BoolNode) ResultTypes() Types {
	return TypeBoolean
}

func (n BoolNode) String() string {
	return "Bool: " + strconv.FormatBool(n.Value)
}

type CeilNode struct {
	Argument Node
}

func (n *CeilNode) ResultTypes() Types {
	return TypeNumber
}

func (n *CeilNode) String() string {
	return "Ceil"
}

func (n *CeilNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type ContainsNode struct {
	Arguments [2]Node
}

func (n *ContainsNode) ResultTypes() Types {
	return TypeBoolean
}

func (n *ContainsNode) String() string {
	return "Contains"
}

func (n *ContainsNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type CurrentNode struct{}

func (n CurrentNode) ResultTypes() Types {
	return TypeAny
}

func (n CurrentNode) String() string {
	return "Current"
}

type DefineVariables struct {
	Variables map[string]Node
	Child     Node
}

func (n *DefineVariables) ResultTypes() Types {
	return n.Child.ResultTypes()
}

func (n *DefineVariables) String() string {
	return "DefineVariables"
}

func (n *DefineVariables) Walk(v Visitor) {
	for _, variable := range n.Variables {
		v.Visit(variable)
	}

	v.Visit(n.Child)
}

type DivideNode struct {
	Left  Node
	Right Node
}

func (n *DivideNode) ResultTypes() Types {
	return TypeNumber
}

func (n *DivideNode) String() string {
	return "Divide"
}

func (n *DivideNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type EndsWithNode struct {
	Arguments [2]Node
}

func (n *EndsWithNode) ResultTypes() Types {
	return TypeBoolean
}

func (n *EndsWithNode) String() string {
	return "EndsWith"
}

func (n *EndsWithNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type EqualNode struct {
	Left  Node
	Right Node
}

func (n *EqualNode) ResultTypes() Types {
	return TypeBoolean
}

func (n *EqualNode) String() string {
	return "Equal"
}

func (n *EqualNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type FieldNode struct {
	Value string
}

func (n *FieldNode) ResultTypes() Types {
	return TypeAny
}

func (n *FieldNode) String() string {
	return "Field: " + n.Value
}

type FilterNode struct {
	Child  Node
	Filter Node
}

func (n *FilterNode) ResultTypes() Types {
	return TypeArray
}

func (n *FilterNode) String() string {
	return "Filter"
}

func (n *FilterNode) Walk(v Visitor) {
	v.Visit(n.Child)
	v.Visit(n.Filter)
}

type FilterAndProjectNode struct {
	Left   Node
	Filter Node
	Right  Node
}

func (n *FilterAndProjectNode) ResultTypes() Types {
	return TypeArray
}

func (n *FilterAndProjectNode) String() string {
	return "FilterAndProject"
}

func (n *FilterAndProjectNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Filter)
	v.Visit(n.Right)
}

type FilterAndProjectCurrentNode struct {
	Filter Node
	Child  Node
}

func (n *FilterAndProjectCurrentNode) ResultTypes() Types {
	return TypeArray
}

func (n *FilterAndProjectCurrentNode) String() string {
	return "FilterAndProjectCurrent"
}

func (n *FilterAndProjectCurrentNode) Walk(v Visitor) {
	v.Visit(n.Filter)
	v.Visit(n.Child)
}

type FilterCurrentNode struct {
	Filter Node
}

func (n *FilterCurrentNode) ResultTypes() Types {
	return TypeArray
}

func (n *FilterCurrentNode) String() string {
	return "FilterCurrent"
}

func (n *FilterCurrentNode) Walk(v Visitor) {
	v.Visit(n.Filter)
}

type FindFirstNode struct {
	Arguments [2]Node
}

func (n *FindFirstNode) ResultTypes() Types {
	return TypeNumber
}

func (n *FindFirstNode) String() string {
	return "FindFirst"
}

func (n *FindFirstNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type FindFirstBetweenNode struct {
	Arguments [4]Node
}

func (n *FindFirstBetweenNode) ResultTypes() Types {
	return TypeNumber
}

func (n *FindFirstBetweenNode) String() string {
	return "FindFirstBetween"
}

func (n *FindFirstBetweenNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
	v.Visit(n.Arguments[2])
	v.Visit(n.Arguments[3])
}

type FindFirstFromNode struct {
	Arguments [3]Node
}

func (n *FindFirstFromNode) ResultTypes() Types {
	return TypeNumber
}

func (n *FindFirstFromNode) String() string {
	return "FindFirstFrom"
}

func (n *FindFirstFromNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
	v.Visit(n.Arguments[2])
}

type FindLastNode struct {
	Arguments [2]Node
}

func (n *FindLastNode) ResultTypes() Types {
	return TypeNumber
}

func (n *FindLastNode) String() string {
	return "FindLast"
}

func (n *FindLastNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type FindLastBetweenNode struct {
	Arguments [4]Node
}

func (n *FindLastBetweenNode) ResultTypes() Types {
	return TypeNumber
}

func (n *FindLastBetweenNode) String() string {
	return "FindLastBetween"
}

func (n *FindLastBetweenNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
	v.Visit(n.Arguments[2])
	v.Visit(n.Arguments[3])
}

type FindLastFromNode struct {
	Arguments [3]Node
}

func (n *FindLastFromNode) ResultTypes() Types {
	return TypeNumber
}

func (n *FindLastFromNode) String() string {
	return "FindLastFrom"
}

func (n *FindLastFromNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
	v.Visit(n.Arguments[2])
}

type FlattenNode struct {
	Child Node
}

func (n *FlattenNode) ResultTypes() Types {
	return TypeArray
}

func (n *FlattenNode) String() string {
	return "Flatten"
}

func (n *FlattenNode) Walk(v Visitor) {
	v.Visit(n.Child)
}

type FlattenAndProjectNode struct {
	Left  Node
	Right Node
}

func (n *FlattenAndProjectNode) ResultTypes() Types {
	return TypeArray
}

func (n *FlattenAndProjectNode) String() string {
	return "FlattenAndProject"
}

func (n *FlattenAndProjectNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type FlattenAndProjectCurrentNode struct {
	Child Node
}

func (n *FlattenAndProjectCurrentNode) ResultTypes() Types {
	return TypeArray
}

func (n *FlattenAndProjectCurrentNode) String() string {
	return "FlattenAndProjectCurrent"
}

func (n *FlattenAndProjectCurrentNode) Walk(v Visitor) {
	v.Visit(n.Child)
}

type FlattenCurrentNode struct{}

func (n FlattenCurrentNode) ResultTypes() Types {
	return TypeArray
}

func (n FlattenCurrentNode) String() string {
	return "FlattenCurrent"
}

type FloorNode struct {
	Argument Node
}

func (n *FloorNode) ResultTypes() Types {
	return TypeNumber
}

func (n *FloorNode) String() string {
	return "Floor"
}

func (n *FloorNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type FromItemsNode struct {
	Argument Node
}

func (n *FromItemsNode) ResultTypes() Types {
	return TypeObject
}

func (n *FromItemsNode) String() string {
	return "FromItems"
}

func (n *FromItemsNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type GreaterNode struct {
	Left  Node
	Right Node
}

func (n *GreaterNode) ResultTypes() Types {
	return TypeBoolean
}

func (n *GreaterNode) String() string {
	return "Greater"
}

func (n *GreaterNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type GreaterOrEqualNode struct {
	Left  Node
	Right Node
}

func (n *GreaterOrEqualNode) ResultTypes() Types {
	return TypeBoolean
}

func (n *GreaterOrEqualNode) String() string {
	return "GreaterOrEqual"
}

func (n *GreaterOrEqualNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type GroupByNode struct {
	Arguments [2]Node
}

func (n *GroupByNode) ResultTypes() Types {
	return TypeObject
}

func (n *GroupByNode) String() string {
	return "GroupBy"
}

func (n *GroupByNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type IfNode struct {
	Condition Node
	Then      Node
	Else      Node
}

func (n *IfNode) String() string {
	return "If"
}

func (n *IfNode) ResultTypes() Types {
	return TypeAny
}

func (n *IfNode) Walk(v Visitor) {
	v.Visit(n.Condition)
	v.Visit(n.Then)
	v.Visit(n.Else)
}

type IndexNode struct {
	Child Node
	Value int
}

func (n *IndexNode) ResultTypes() Types {
	return TypeAny
}

func (n *IndexNode) String() string {
	return "Index: " + strconv.FormatInt(int64(n.Value), 10)
}

func (n *IndexNode) Walk(v Visitor) {
	v.Visit(n.Child)
}

type IndexCurrentNode struct {
	Value int
}

func (n IndexCurrentNode) ResultTypes() Types {
	return TypeAny
}

func (n *IndexCurrentNode) String() string {
	return "IndexCurrent: " + strconv.FormatInt(int64(n.Value), 10)
}

type IntegerDivideNode struct {
	Left  Node
	Right Node
}

func (n *IntegerDivideNode) ResultTypes() Types {
	return TypeNumber
}

func (n *IntegerDivideNode) String() string {
	return "IntegerDivide"
}

func (n *IntegerDivideNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type ItemsNode struct {
	Argument Node
}

func (n *ItemsNode) ResultTypes() Types {
	return TypeArray
}

func (n *ItemsNode) String() string {
	return "Items"
}

func (n *ItemsNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type JoinNode struct {
	Arguments [2]Node
}

func (n *JoinNode) ResultTypes() Types {
	return TypeString
}

func (n *JoinNode) String() string {
	return "Join"
}

func (n *JoinNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type KeysNode struct {
	Argument Node
}

func (n *KeysNode) ResultTypes() Types {
	return TypeArray
}

func (n *KeysNode) String() string {
	return "Keys"
}

func (n *KeysNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type LengthNode struct {
	Argument Node
}

func (n *LengthNode) ResultTypes() Types {
	return TypeNumber
}

func (n *LengthNode) String() string {
	return "Length"
}

func (n *LengthNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type LessNode struct {
	Left  Node
	Right Node
}

func (n *LessNode) ResultTypes() Types {
	return TypeBoolean
}

func (n *LessNode) String() string {
	return "Less"
}

func (n *LessNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type LessOrEqualNode struct {
	Left  Node
	Right Node
}

func (n *LessOrEqualNode) ResultTypes() Types {
	return TypeBoolean
}

func (n *LessOrEqualNode) String() string {
	return "LessOrEqual"
}

func (n *LessOrEqualNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type LowerNode struct {
	Argument Node
}

func (n *LowerNode) ResultTypes() Types {
	return TypeString
}

func (n *LowerNode) String() string {
	return "Lower"
}

func (n *LowerNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type MapNode struct {
	Arguments [2]Node
}

func (n *MapNode) ResultTypes() Types {
	return TypeArray
}

func (n *MapNode) String() string {
	return "Map"
}

func (n *MapNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type MaxNode struct {
	Argument Node
}

func (n *MaxNode) ResultTypes() Types {
	return TypeNumber | TypeString
}

func (n *MaxNode) String() string {
	return "Max"
}

func (n *MaxNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type MaxByNode struct {
	Arguments [2]Node
}

func (n *MaxByNode) ResultTypes() Types {
	return TypeAny
}

func (n *MaxByNode) String() string {
	return "MaxBy"
}

func (n *MaxByNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type MergeNode struct {
	Arguments []Node
}

func (n *MergeNode) ResultTypes() Types {
	return TypeObject
}

func (n *MergeNode) String() string {
	return "Merge"
}

func (n *MergeNode) Walk(v Visitor) {
	for _, arg := range n.Arguments {
		v.Visit(arg)
	}
}

type MinNode struct {
	Argument Node
}

func (n *MinNode) ResultTypes() Types {
	return TypeNumber | TypeString
}

func (n *MinNode) String() string {
	return "Min"
}

func (n *MinNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type MinByNode struct {
	Arguments [2]Node
}

func (n *MinByNode) ResultTypes() Types {
	return TypeAny
}

func (n *MinByNode) String() string {
	return "MinBy"
}

func (n *MinByNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type ModuloNode struct {
	Left  Node
	Right Node
}

func (n *ModuloNode) ResultTypes() Types {
	return TypeNumber
}

func (n *ModuloNode) String() string {
	return "Modulo"
}

func (n *ModuloNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type MultiplyNode struct {
	Left  Node
	Right Node
}

func (n *MultiplyNode) ResultTypes() Types {
	return TypeNumber
}

func (n *MultiplyNode) String() string {
	return "Multiply"
}

func (n *MultiplyNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type NegateNode struct {
	Child Node
}

func (n *NegateNode) ResultTypes() Types {
	return TypeNumber
}

func (n *NegateNode) String() string {
	return "Negate"
}

func (n *NegateNode) Walk(v Visitor) {
	v.Visit(n.Child)
}

type NotNode struct {
	Child Node
}

func (n *NotNode) ResultTypes() Types {
	return TypeBoolean
}

func (n *NotNode) String() string {
	return "Not"
}

func (n *NotNode) Walk(v Visitor) {
	v.Visit(n.Child)
}

type NotEqualNode struct {
	Left  Node
	Right Node
}

func (n *NotEqualNode) ResultTypes() Types {
	return TypeBoolean
}

func (n *NotEqualNode) String() string {
	return "NotEqual"
}

func (n *NotEqualNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type NotNullNode struct {
	Arguments []Node
}

func (n *NotNullNode) ResultTypes() Types {
	return TypeAny
}

func (n *NotNullNode) String() string {
	return "NotNull"
}

func (n *NotNullNode) Walk(v Visitor) {
	for _, arg := range n.Arguments {
		v.Visit(arg)
	}
}

type NotNullValueNode struct {
	Argument Node
	Value    any
}

func (n *NotNullValueNode) ResultTypes() Types {
	return TypeAny
}

func (n *NotNullValueNode) String() string {
	switch value := n.Value.(type) {
	case json.Number:
		return "NotNullValue: " + value.String()
	case string:
		return "NotNullValue: " + value
	default:
		return "NotNullValue"
	}
}

type NullNode struct{}

func (n NullNode) ResultTypes() Types {
	return TypeNull
}

func (n NullNode) String() string {
	return "Null"
}

type ObjectValuesNode struct {
	Child Node
}

func (n *ObjectValuesNode) ResultTypes() Types {
	return TypeArray
}

func (n *ObjectValuesNode) String() string {
	return "ObjectValues"
}

func (n *ObjectValuesNode) Walk(v Visitor) {
	v.Visit(n.Child)
}

type ObjectValuesCurrentNode struct{}

func (n ObjectValuesCurrentNode) ResultTypes() Types {
	return TypeArray
}

func (n ObjectValuesCurrentNode) String() string {
	return "ObjectValuesCurrent"
}

type OrNode struct {
	Left  Node
	Right Node
}

func (n *OrNode) ResultTypes() Types {
	return TypeAny
}

func (n *OrNode) String() string {
	return "Or"
}

func (n *OrNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type PadLeftNode struct {
	Arguments [3]Node
}

func (n *PadLeftNode) ResultTypes() Types {
	return TypeString
}

func (n *PadLeftNode) String() string {
	return "PadLeft"
}

func (n *PadLeftNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
	v.Visit(n.Arguments[2])
}

type PadRightNode struct {
	Arguments [3]Node
}

func (n *PadRightNode) ResultTypes() Types {
	return TypeString
}

func (n *PadRightNode) String() string {
	return "PadRight"
}

func (n *PadRightNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
	v.Visit(n.Arguments[2])
}

type PadSpaceLeftNode struct {
	Arguments [2]Node
}

func (n *PadSpaceLeftNode) ResultTypes() Types {
	return TypeString
}

func (n *PadSpaceLeftNode) String() string {
	return "PadSpaceLeft"
}

func (n *PadSpaceLeftNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type PadSpaceRightNode struct {
	Arguments [2]Node
}

func (n *PadSpaceRightNode) ResultTypes() Types {
	return TypeString
}

func (n *PadSpaceRightNode) String() string {
	return "PadSpaceRight"
}

func (n *PadSpaceRightNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type PipeNode struct {
	Left  Node
	Right Node
}

func (n *PipeNode) ResultTypes() Types {
	return TypeAny
}

func (n *PipeNode) String() string {
	return "Pipe"
}

func (n *PipeNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type PipeFieldNode struct {
	Left  Node
	Right string
}

func (n *PipeFieldNode) ResultTypes() Types {
	return TypeAny
}

func (n *PipeFieldNode) String() string {
	return "PipeField: " + n.Right
}

func (n *PipeFieldNode) Walk(v Visitor) {
	v.Visit(n.Left)
}

type ProjectArrayNode struct {
	Left  Node
	Right Node
}

func (n *ProjectArrayNode) ResultTypes() Types {
	switch n.Left.(type) {
	case *SliceNode,
		*SliceCurrentNode,
		*SliceStepNode,
		*SliceStepCurrentNode:
		return TypeAny
	}

	return TypeArray
}

func (n *ProjectArrayNode) String() string {
	return "ProjectArray"
}

func (n *ProjectArrayNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type ProjectArrayCurrentNode struct {
	Child Node
}

func (n *ProjectArrayCurrentNode) ResultTypes() Types {
	return TypeArray
}

func (n *ProjectArrayCurrentNode) String() string {
	return "ProjectArrayCurrent"
}

func (n *ProjectArrayCurrentNode) Walk(v Visitor) {
	v.Visit(n.Child)
}

type ProjectObjectNode struct {
	Left  Node
	Right Node
}

func (n *ProjectObjectNode) ResultTypes() Types {
	return TypeArray
}

func (n *ProjectObjectNode) String() string {
	return "ProjectObject"
}

func (n *ProjectObjectNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type ProjectObjectCurrentNode struct {
	Child Node
}

func (n *ProjectObjectCurrentNode) ResultTypes() Types {
	return TypeArray
}

func (n *ProjectObjectCurrentNode) String() string {
	return "ProjectObjectCurrent"
}

func (n *ProjectObjectCurrentNode) Walk(v Visitor) {
	v.Visit(n.Child)
}

type PruneArrayNode struct {
	Child Node
}

func (n *PruneArrayNode) ResultTypes() Types {
	return TypeArray
}

func (n *PruneArrayNode) String() string {
	return "PruneArray"
}

func (n *PruneArrayNode) Walk(v Visitor) {
	v.Visit(n.Child)
}

type PruneArrayCurrentNode struct{}

func (n PruneArrayCurrentNode) ResultTypes() Types {
	return TypeArray
}

func (n PruneArrayCurrentNode) String() string {
	return "PruneArrayCurrent"
}

type ReplaceNode struct {
	Arguments [3]Node
}

func (n *ReplaceNode) ResultTypes() Types {
	return TypeString
}

func (n *ReplaceNode) String() string {
	return "Replace"
}

func (n *ReplaceNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
	v.Visit(n.Arguments[2])
}

type ReplaceCountNode struct {
	Arguments [4]Node
}

func (n *ReplaceCountNode) ResultTypes() Types {
	return TypeString
}

func (n *ReplaceCountNode) String() string {
	return "ReplaceCount"
}

func (n *ReplaceCountNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
	v.Visit(n.Arguments[2])
	v.Visit(n.Arguments[3])
}

type ReverseNode struct {
	Argument Node
}

func (n *ReverseNode) ResultTypes() Types {
	return TypeString | TypeArray
}

func (n *ReverseNode) String() string {
	return "Reverse"
}

func (n *ReverseNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type RootNode struct{}

func (n RootNode) ResultTypes() Types {
	return TypeAny
}

func (n RootNode) String() string {
	return "Root"
}

type SelectArrayNode struct {
	Child  Node
	Fields []Node
}

func (n *SelectArrayNode) ResultTypes() Types {
	return TypeArray
}

func (n *SelectArrayNode) String() string {
	return "SelectArray"
}

func (n *SelectArrayNode) Walk(v Visitor) {
	v.Visit(n.Child)

	for _, field := range n.Fields {
		v.Visit(field)
	}
}

type SelectArrayCurrentNode struct {
	Fields []Node
}

func (n *SelectArrayCurrentNode) ResultTypes() Types {
	return TypeArray
}

func (n *SelectArrayCurrentNode) String() string {
	return "SelectArrayCurrent"
}

func (n *SelectArrayCurrentNode) Walk(v Visitor) {
	for _, field := range n.Fields {
		v.Visit(field)
	}
}

type SelectArraySingleNode struct {
	Child Node
	Field Node
}

func (n *SelectArraySingleNode) ResultTypes() Types {
	return TypeArray
}

func (n *SelectArraySingleNode) String() string {
	return "SelectArraySingle"
}

func (n *SelectArraySingleNode) Walk(v Visitor) {
	v.Visit(n.Child)
	v.Visit(n.Field)
}

type SelectArraySingleCurrentNode struct {
	Field Node
}

func (n *SelectArraySingleCurrentNode) ResultTypes() Types {
	return TypeArray
}

func (n *SelectArraySingleCurrentNode) String() string {
	return "SelectArraySingleCurrent"
}

func (n *SelectArraySingleCurrentNode) Walk(v Visitor) {
	v.Visit(n.Field)
}

type SelectObjectNode struct {
	Child  Node
	Fields map[string]Node
}

func (n *SelectObjectNode) ResultTypes() Types {
	return TypeObject
}

func (n *SelectObjectNode) String() string {
	return "SelectObject"
}

func (n *SelectObjectNode) Walk(v Visitor) {
	v.Visit(n.Child)

	for _, field := range n.Fields {
		v.Visit(field)
	}
}

type SelectObjectCurrentNode struct {
	Fields map[string]Node
}

func (n *SelectObjectCurrentNode) ResultTypes() Types {
	return TypeObject
}

func (n *SelectObjectCurrentNode) String() string {
	return "SelectObjectCurrent"
}

func (n *SelectObjectCurrentNode) Walk(v Visitor) {
	for _, field := range n.Fields {
		v.Visit(field)
	}
}

type SelectObjectSingleNode struct {
	Child Node
	Key   string
	Field Node
}

func (n *SelectObjectSingleNode) ResultTypes() Types {
	return TypeObject
}

func (n *SelectObjectSingleNode) String() string {
	return "SelectObjectSingle"
}

func (n *SelectObjectSingleNode) Walk(v Visitor) {
	v.Visit(n.Child)
	v.Visit(n.Field)
}

type SelectObjectSingleCurrentNode struct {
	Key   string
	Field Node
}

func (n *SelectObjectSingleCurrentNode) ResultTypes() Types {
	return TypeObject
}

func (n *SelectObjectSingleCurrentNode) String() string {
	return "SelectObjectSingleCurrent"
}

func (n *SelectObjectSingleCurrentNode) Walk(v Visitor) {
	v.Visit(n.Field)
}

type SliceNode struct {
	Child Node
	Start int
	Stop  int
}

func (n *SliceNode) ResultTypes() Types {
	return TypeString | TypeArray
}

func (n *SliceNode) String() string {
	return "Slice: " + strconv.Itoa(n.Start) + ":" + strconv.Itoa(n.Stop)
}

func (n *SliceNode) Walk(v Visitor) {
	v.Visit(n.Child)
}

type SliceCurrentNode struct {
	Start int
	Stop  int
}

func (n *SliceCurrentNode) ResultTypes() Types {
	return TypeString | TypeArray
}

func (n *SliceCurrentNode) String() string {
	return "SliceCurrent: " + strconv.Itoa(n.Start) + ":" + strconv.Itoa(n.Stop)
}

type SliceStepNode struct {
	Child Node
	Start int
	Stop  int
	Step  int
}

func (n *SliceStepNode) ResultTypes() Types {
	return TypeString | TypeArray
}

func (n *SliceStepNode) String() string {
	return "SliceStep: " + strconv.Itoa(n.Start) + ":" + strconv.Itoa(n.Stop) + ":" + strconv.Itoa(n.Step)
}

func (n *SliceStepNode) Walk(v Visitor) {
	v.Visit(n.Child)
}

type SliceStepCurrentNode struct {
	Start int
	Stop  int
	Step  int
}

func (n *SliceStepCurrentNode) ResultTypes() Types {
	return TypeString | TypeArray
}

func (n *SliceStepCurrentNode) String() string {
	return "SliceStepCurrent: " + strconv.Itoa(n.Start) + ":" + strconv.Itoa(n.Stop) + ":" + strconv.Itoa(n.Step)
}

type SmallIndexCurrentNode struct {
	Value uint8
}

func (n SmallIndexCurrentNode) ResultTypes() Types {
	return TypeAny
}

func (n SmallIndexCurrentNode) String() string {
	return "SmallIndexCurrent: " + strconv.FormatUint(uint64(n.Value), 10)
}

type SortNode struct {
	Argument Node
}

func (n *SortNode) ResultTypes() Types {
	return TypeArray
}

func (n *SortNode) String() string {
	return "Sort"
}

func (n *SortNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type SortByNode struct {
	Arguments [2]Node
}

func (n *SortByNode) ResultTypes() Types {
	return TypeArray
}

func (n *SortByNode) String() string {
	return "SortBy"
}

func (n *SortByNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type SplitNode struct {
	Arguments [2]Node
}

func (n *SplitNode) ResultTypes() Types {
	return TypeArray
}

func (n *SplitNode) String() string {
	return "Split"
}

func (n *SplitNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type SplitCountNode struct {
	Arguments [3]Node
}

func (n *SplitCountNode) ResultTypes() Types {
	return TypeArray
}

func (n *SplitCountNode) String() string {
	return "SplitCount"
}

func (n *SplitCountNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
	v.Visit(n.Arguments[2])
}

type StartsWithNode struct {
	Arguments [2]Node
}

func (n *StartsWithNode) ResultTypes() Types {
	return TypeBoolean
}

func (n *StartsWithNode) String() string {
	return "StartsWith"
}

func (n *StartsWithNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type SubtractNode struct {
	Left  Node
	Right Node
}

func (n *SubtractNode) ResultTypes() Types {
	return TypeNumber
}

func (n *SubtractNode) String() string {
	return "Subtract"
}

func (n *SubtractNode) Walk(v Visitor) {
	v.Visit(n.Left)
	v.Visit(n.Right)
}

type SumNode struct {
	Argument Node
}

func (n *SumNode) ResultTypes() Types {
	return TypeNumber
}

func (n *SumNode) String() string {
	return "Sum"
}

func (n *SumNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type ToArrayNode struct {
	Argument Node
}

func (n *ToArrayNode) ResultTypes() Types {
	return TypeArray
}

func (n *ToArrayNode) String() string {
	return "ToArray"
}

func (n *ToArrayNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type ToNumberNode struct {
	Argument Node
}

func (n *ToNumberNode) ResultTypes() Types {
	return TypeNumber
}

func (n *ToNumberNode) String() string {
	return "ToNumber"
}

func (n *ToNumberNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type ToStringNode struct {
	Argument Node
}

func (n *ToStringNode) ResultTypes() Types {
	return TypeString
}

func (n *ToStringNode) String() string {
	return "ToString"
}

func (n *ToStringNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type TrimNode struct {
	Arguments [2]Node
}

func (n *TrimNode) ResultTypes() Types {
	return TypeString
}

func (n *TrimNode) String() string {
	return "Trim"
}

func (n *TrimNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type TrimLeftNode struct {
	Arguments [2]Node
}

func (n *TrimLeftNode) ResultTypes() Types {
	return TypeString
}

func (n *TrimLeftNode) String() string {
	return "TrimLeft"
}

func (n *TrimLeftNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type TrimRightNode struct {
	Arguments [2]Node
}

func (n *TrimRightNode) ResultTypes() Types {
	return TypeString
}

func (n *TrimRightNode) String() string {
	return "TrimRight"
}

func (n *TrimRightNode) Walk(v Visitor) {
	v.Visit(n.Arguments[0])
	v.Visit(n.Arguments[1])
}

type TrimSpaceNode struct {
	Argument Node
}

func (n *TrimSpaceNode) ResultTypes() Types {
	return TypeString
}

func (n *TrimSpaceNode) String() string {
	return "TrimSpace"
}

func (n *TrimSpaceNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type TrimSpaceLeftNode struct {
	Argument Node
}

func (n *TrimSpaceLeftNode) ResultTypes() Types {
	return TypeString
}

func (n *TrimSpaceLeftNode) String() string {
	return "TrimSpaceLeft"
}

func (n *TrimSpaceLeftNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type TrimSpaceRightNode struct {
	Argument Node
}

func (n *TrimSpaceRightNode) ResultTypes() Types {
	return TypeString
}

func (n *TrimSpaceRightNode) String() string {
	return "TrimSpaceRight"
}

func (n *TrimSpaceRightNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type TypeNode struct {
	Argument Node
}

func (n *TypeNode) ResultTypes() Types {
	return TypeString
}

func (n *TypeNode) String() string {
	return "Type"
}

func (n *TypeNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type UpperNode struct {
	Argument Node
}

func (n *UpperNode) ResultTypes() Types {
	return TypeString
}

func (n *UpperNode) String() string {
	return "Upper"
}

func (n *UpperNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type ValueNode struct {
	Value any
}

func (n *ValueNode) ResultTypes() Types {
	switch n.Value.(type) {
	case []any:
		return TypeArray
	case map[string]any:
		return TypeObject
	case json.Number:
		return TypeNumber
	case string:
		return TypeString
	default:
		return TypeAny
	}
}

func (n *ValueNode) String() string {
	switch value := n.Value.(type) {
	case json.Number:
		return "Value: " + value.String()
	case string:
		return "Value: " + value
	default:
		return "Value"
	}
}

type ValuesNode struct {
	Argument Node
}

func (n *ValuesNode) ResultTypes() Types {
	return TypeArray
}

func (n *ValuesNode) String() string {
	return "Values"
}

func (n *ValuesNode) Walk(v Visitor) {
	v.Visit(n.Argument)
}

type VariableNode struct {
	Name string
}

func (n *VariableNode) ResultTypes() Types {
	return TypeAny
}

func (n *VariableNode) String() string {
	return "Variable: " + n.Name
}

type ZipNode struct {
	Arguments []Node
}

func (n *ZipNode) ResultTypes() Types {
	return TypeArray
}

func (n *ZipNode) String() string {
	return "Zip"
}

func (n *ZipNode) Walk(v Visitor) {
	for _, arg := range n.Arguments {
		v.Visit(arg)
	}
}
