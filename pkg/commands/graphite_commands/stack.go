package graphite_commands

// Stack contains methods for interacting with Graphite stacks
type Stack struct {
	*Common
}

// NewStack creates a new Stack instance
func NewStack(common *Common) *Stack {
	return &Stack{
		Common: common,
	}
}

// Submit submits a stack to the remote repository
func (s *Stack) Submit() error {
	cmdArgs := NewGraphiteCmd("submit").
		Arg("--stack").
		ToArgv()

	return s.GetCmd().New(cmdArgs).Run()
}

// Restack ensures all branches in a stack are based on the current version of their parents
func (s *Stack) Restack() error {
	cmdArgs := NewGraphiteCmd("restack").
		ToArgv()

	return s.GetCmd().New(cmdArgs).Run()
}

// Log shows the commit history for the current stack
func (s *Stack) Log(short bool) (string, error) {
	cmdArgs := NewGraphiteCmd("log").
		ArgIf(short, "short").
		ArgIf(!short, "long").
		ToArgv()

	return s.GetCmd().New(cmdArgs).RunWithOutput()
}
