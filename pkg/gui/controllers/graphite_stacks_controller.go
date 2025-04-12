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
	}
}

func (self *GraphiteStacksController) checkout(stack *models.GraphiteStack) error {
	// Use the RefsHelper to checkout the branch
	return self.c.Helpers().Refs.CheckoutRef(stack.Name, types.CheckoutRefOptions{})
}
