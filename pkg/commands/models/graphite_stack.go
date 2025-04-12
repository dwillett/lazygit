package models

// GraphiteStack represents a stack of branches in Graphite
type GraphiteStack struct {
	// Name of the branch
	Name string
	// Whether this is the current branch
	IsCurrent bool
	// Whether this is a trunk branch
	IsTrunk bool
	// Parent branches in the stack
	Parents []*GraphiteStack
	// Child branches in the stack
	Children []*GraphiteStack
}

func (s *GraphiteStack) ID() string {
	return s.Name
}

func (s *GraphiteStack) Description() string {
	return s.Name
}

func (s GraphiteStack) URN() string {
	return s.Name
}
