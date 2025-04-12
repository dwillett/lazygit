package graphite_commands

import "github.com/jesseduffield/lazygit/pkg/commands/oscommands"

// Branch contains methods for interacting with Graphite branches
type Branch struct {
	*Common
}

// NewBranch creates a new Branch instance
func NewBranch(common *Common) *Branch {
	return &Branch{
		Common: common,
	}
}

// Checkout checks out a branch
func (b *Branch) Checkout(branchName string) error {
	cmdArgs := NewGraphiteCmd("checkout").
		Arg(branchName).
		ToArgv()

	return b.GetCmd().New(cmdArgs).Run()
}

// Create creates a new branch with a commit
func (b *Branch) Create(message string, all bool) error {
	cmdArgs := NewGraphiteCmd("create").
		ArgIf(all, "--all").
		Arg("-m", message).
		ToArgv()

	return b.GetCmd().New(cmdArgs).Run()
}

// Info gets information about a branch
func (b *Branch) Info(branchName string) (string, error) {
	cmdArgs := NewGraphiteCmd("info").
		Arg(branchName).
		ToArgv()

	return b.GetCmd().New(cmdArgs).RunWithOutput()
}

// Track tracks an existing git branch as a Graphite branch
func (b *Branch) Track(branchName string, parentBranchName string) error {
	cmdArgs := NewGraphiteCmd("track").
		Arg(branchName).
		ArgIf(parentBranchName != "", "--parent", parentBranchName).
		ToArgv()

	return b.GetCmd().New(cmdArgs).Run()
}

// Untrack removes Graphite metadata for a branch
func (b *Branch) Untrack(branchName string) error {
	cmdArgs := NewGraphiteCmd("untrack").
		Arg(branchName).
		ToArgv()

	return b.GetCmd().New(cmdArgs).Run()
}

// Rename renames a branch
func (b *Branch) Rename(oldName, newName string) error {
	cmdArgs := NewGraphiteCmd("rename").
		Arg(oldName, newName).
		ToArgv()

	return b.GetCmd().New(cmdArgs).Run()
}

// Modify modifies the current branch's commit
func (b *Branch) Modify() error {
	return b.ModifyCmdObj().Run()
}

func (b *Branch) ModifyCmdObj() oscommands.ICmdObj {
	cmdArgs := NewGraphiteCmd("modify").
		ToArgv()

	return b.GetCmd().New(cmdArgs)
}
