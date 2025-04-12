package graphite_commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type GraphiteStacksLoader struct {
	*common.Common
	cmd       oscommands.ICmdObjBuilder
	gitCommon *git_commands.GitCommon
}

func NewGraphiteStacksLoader(common *common.Common, cmd oscommands.ICmdObjBuilder, gitCommon *git_commands.GitCommon) *GraphiteStacksLoader {
	return &GraphiteStacksLoader{
		Common:    common,
		cmd:       cmd,
		gitCommon: gitCommon,
	}
}

// GetStacks returns all Graphite stacks
func (l *GraphiteStacksLoader) GetStacks() ([]*models.GraphiteStack, error) {
	// Create a Graphite Common instance
	graphiteCommon := New(l.Common, l.cmd, l.gitCommon)

	// Create a Repo instance
	repo := NewRepo(graphiteCommon)

	// Use the State method to get the stack structure with trunk at
	stacks, err := repo.State()
	if err != nil {
		l.Log.Error(err)
		return nil, err
	}

	return stacks, nil
}
