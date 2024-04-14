package parser

func isProjectNode(node Node) bool {
	switch node.(type) {
	case *FilterAndProjectNode,
		*FilterAndProjectCurrentNode,
		*FlattenAndProjectNode,
		*FlattenAndProjectCurrentNode,
		*ProjectArrayNode,
		*ProjectArrayCurrentNode:
		return true
	}

	return false
}
