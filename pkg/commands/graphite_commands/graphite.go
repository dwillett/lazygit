package graphite_commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

// GraphiteCommand is our main Graphite interface
type GraphiteCommand struct {
	Branch          *Branch
	Stack           *Stack
	Repo            *Repo
	StackNavigation *StackNavigation
}

// NewGraphiteCommand creates a new GraphiteCommand instance
func NewGraphiteCommand(
	cmn *common.Common,
	cmd oscommands.ICmdObjBuilder,
	gitCommon *git_commands.GitCommon,
) *GraphiteCommand {
	// Create a Graphite Common instance
	graphiteCommon := New(cmn, cmd, gitCommon)

	// Create instances of each command type
	branch := NewBranch(graphiteCommon)
	stack := NewStack(graphiteCommon)
	repo := NewRepo(graphiteCommon)
	stackNavigation := NewStackNavigation(graphiteCommon)

	return &GraphiteCommand{
		Branch:          branch,
		Stack:           stack,
		Repo:            repo,
		StackNavigation: stackNavigation,
	}
}
