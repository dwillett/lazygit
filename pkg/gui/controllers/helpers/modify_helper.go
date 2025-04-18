package helpers

import "github.com/jesseduffield/lazygit/pkg/commands/git_commands"

type ModifyHelper struct {
	c   *HelperCommon
	gpg *GpgHelper
}

func NewModifyHelper(
	c *HelperCommon,
	gpg *GpgHelper,
) *ModifyHelper {
	return &ModifyHelper{
		c:   c,
		gpg: gpg,
	}
}

func (self *ModifyHelper) ModifyBranch() error {
	cmdObj := self.c.Git().Graphite.Branch.ModifyCmdObj()
	self.c.LogAction(self.c.Tr.Actions.ModifyBranch)
	return self.gpg.WithGpgHandling(cmdObj, git_commands.CommitGpgSign, self.c.Tr.AmendingStatus, nil, nil)
}
