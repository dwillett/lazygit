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

// NewGraphiteCommand returns a new GraphiteCommand
func NewGraphiteCommand(cmn *common.Common, cmd oscommands.ICmdObjBuilder, gitCommon *git_commands.GitCommon) *GraphiteCommand {
	graphiteCommon := New(cmn, cmd, gitCommon)
	return &GraphiteCommand{
		Branch:          NewBranch(graphiteCommon),
		Stack:           NewStack(graphiteCommon),
		Repo:            NewRepo(graphiteCommon),
		StackNavigation: NewStackNavigation(graphiteCommon),
	}
}
