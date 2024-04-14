package evaluator

type variableScope struct {
	parent    *variableScope
	variables map[string]any
}

func (s *variableScope) get(variable string) (any, bool) {
	if s == nil {
		return nil, false
	}

	if value, ok := s.variables[variable]; ok {
		return value, true
	}

	if s.parent != nil {
		return s.parent.get(variable)
	}

	return nil, false
}

func (s *variableScope) new(variables map[string]any) *variableScope {
	return &variableScope{
		parent:    s,
		variables: variables,
	}
}
