package graphite_commands

// StackNavigation contains methods for navigating between branches in a stack
type StackNavigation struct {
	*Common
}

// NewStackNavigation creates a new StackNavigation instance
func NewStackNavigation(common *Common) *StackNavigation {
	return &StackNavigation{
		Common: common,
	}
}

// Up moves up one branch in the stack
func (sn *StackNavigation) Up() error {
	cmdArgs := NewGraphiteCmd("up").
		ToArgv()

	return sn.GetCmd().New(cmdArgs).Run()
}

// Down moves down one branch in the stack
func (sn *StackNavigation) Down() error {
	cmdArgs := NewGraphiteCmd("down").
		ToArgv()

	return sn.GetCmd().New(cmdArgs).Run()
}

// UpstackOnto moves the current branch and all branches above it onto a new base
func (sn *StackNavigation) UpstackOnto(baseBranch string) error {
	cmdArgs := NewGraphiteCmd("upstack").
		Arg("onto", baseBranch).
		ToArgv()

	return sn.GetCmd().New(cmdArgs).Run()
}

// DownstackGet pulls in new changes from the remote repository for the current branch and all branches below it
func (sn *StackNavigation) DownstackGet() error {
	cmdArgs := NewGraphiteCmd("downstack").
		Arg("get").
		ToArgv()

	return sn.GetCmd().New(cmdArgs).Run()
}

// DownstackEdit opens an editor to edit the commit message for the current branch and all branches below it
func (sn *StackNavigation) DownstackEdit() error {
	cmdArgs := NewGraphiteCmd("downstack").
		Arg("edit").
		ToArgv()

	return sn.GetCmd().New(cmdArgs).Run()
}
