package graphite_commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

// Common contains dependencies that are shared across all Graphite commands
type Common struct {
	*common.Common
	cmd       oscommands.ICmdObjBuilder
	GitCommon *git_commands.GitCommon
}

// New creates a new Common instance
func New(common *common.Common, cmd oscommands.ICmdObjBuilder, gitCommon *git_commands.GitCommon) *Common {
	return &Common{
		Common:    common,
		cmd:       cmd,
		GitCommon: gitCommon,
	}
}

// GetCmd returns the command builder
func (c *Common) GetCmd() oscommands.ICmdObjBuilder {
	return c.cmd
}
