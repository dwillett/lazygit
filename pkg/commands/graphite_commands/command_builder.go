package graphite_commands

import (
	"strings"
)

// GraphiteCommandBuilder is a convenience struct for building Graphite commands
type GraphiteCommandBuilder struct {
	args []string
}

// NewGraphiteCmd creates a new GraphiteCommandBuilder with the given command
func NewGraphiteCmd(command string) *GraphiteCommandBuilder {
	return &GraphiteCommandBuilder{args: []string{command, "--quiet"}}
}

// Arg adds arguments to the command
func (self *GraphiteCommandBuilder) Arg(args ...string) *GraphiteCommandBuilder {
	self.args = append(self.args, args...)
	return self
}

// ArgIf adds arguments conditionally
func (self *GraphiteCommandBuilder) ArgIf(condition bool, ifTrue ...string) *GraphiteCommandBuilder {
	if condition {
		self.Arg(ifTrue...)
	}
	return self
}

// ArgIfElse adds one of two arguments based on a condition
func (self *GraphiteCommandBuilder) ArgIfElse(condition bool, ifTrue string, ifFalse string) *GraphiteCommandBuilder {
	if condition {
		return self.Arg(ifTrue)
	} else {
		return self.Arg(ifFalse)
	}
}

// ToArgv returns the command as a slice of strings
func (self *GraphiteCommandBuilder) ToArgv() []string {
	return append([]string{"gt"}, self.args...)
}

// ToString returns the command as a string
func (self *GraphiteCommandBuilder) ToString() string {
	return strings.Join(self.ToArgv(), " ")
}
