package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type GraphiteStacksContext struct {
	*FilteredListViewModel[*models.GraphiteStack]
	*ListContextTrait
}

var _ types.IListContext = (*GraphiteStacksContext)(nil)

func NewGraphiteStacksContext(c *ContextCommon) *GraphiteStacksContext {
	viewModel := NewFilteredListViewModel(
		func() []*models.GraphiteStack { return c.Model().GraphiteStacks },
		func(stack *models.GraphiteStack) []string {
			return []string{stack.Name}
		},
	)

	getDisplayStrings := func(_ int, _ int) [][]string {
		return presentation.GetGraphiteStackListDisplayStrings(
			viewModel.GetItems(),
			c.State().GetItemOperation,
			c.State().GetRepoState().GetScreenMode() != types.SCREEN_NORMAL,
			c.Views().GraphiteStacks.InnerWidth(),
			c.Tr,
			c.UserConfig(),
		)
	}
	self := &GraphiteStacksContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:                       c.Views().GraphiteStacks,
				WindowName:                 "graphiteStacks",
				Key:                        GRAPHITE_STACKS_CONTEXT_KEY,
				Kind:                       types.SIDE_CONTEXT,
				Focusable:                  true,
				NeedsRerenderOnWidthChange: types.NEEDS_RERENDER_ON_WIDTH_CHANGE_WHEN_WIDTH_CHANGES,
			})),
			ListRenderer: ListRenderer{
				list:              viewModel,
				getDisplayStrings: getDisplayStrings,
			},
			c: c,
		},
	}

	return self
}
