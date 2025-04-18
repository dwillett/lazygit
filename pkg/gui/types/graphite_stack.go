package types

// GraphiteStack represents a stack of branches in Graphite
type GraphiteStack interface {
	ID() string
	Name() string
	Branches() []string
	URN() string
}
