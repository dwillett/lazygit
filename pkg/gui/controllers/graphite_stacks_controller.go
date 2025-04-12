package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type GraphiteStacksController struct {
	baseController
	*ListControllerTrait[*models.GraphiteStack]
	c *ControllerCommon
}

func NewGraphiteStacksController(c *ControllerCommon) *GraphiteStacksController {
	return &GraphiteStacksController{
		baseController: baseController{},
		c:              c,
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().GraphiteStacks,
			c.Contexts().GraphiteStacks.GetSelected,
			c.Contexts().GraphiteStacks.GetSelectedItems,
		),
	}
}

func (self *GraphiteStacksController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:               opts.GetKey(opts.Config.Universal.Select),
			Handler:           self.ListControllerTrait.withItem(self.checkout),
			GetDisabledReason: self.ListControllerTrait.require(self.ListControllerTrait.singleItemSelected()),
			Description:       self.c.Tr.Checkout,
			Tooltip:           self.c.Tr.CheckoutTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Graphite.Modify),
			Handler:           self.withItem(self.modify),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Modify,
			Tooltip:           self.c.Tr.ModifyCommitTooltip,
			DisplayOnScreen:   true,
		},
	}
}

func (self *GraphiteStacksController) checkout(stack *models.GraphiteStack) error {
	// Use the RefsHelper to checkout the branch
	return self.c.Helpers().Refs.CheckoutRef(stack.Name, types.CheckoutRefOptions{})
}

func (self *GraphiteStacksController) modify(stack *models.GraphiteStack) error {
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.ModifyCommitTitle,
		Prompt: self.c.Tr.ModifyCommitPrompt,
		HandleConfirm: func() error {
			return self.c.Helpers().WorkingTree.WithEnsureCommittableFiles(func() error {
				if err := self.c.Helpers().ModifyHelper.ModifyBranch(); err != nil {
					return err
				}
				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
			})
		},
	})

	return nil
}
