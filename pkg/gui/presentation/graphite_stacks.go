package presentation

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

const (
	// Tree visualization characters
	stackIndicatorEmpty   = "◯"
	stackIndicatorCurrent = "◉"
	stackVerticalLine     = "│ "
	stackHorizontalLine   = "─┴"
	stackEndLine          = "─┘"
)

// stackColors defines the colors to cycle through for different stack levels
var stackColors = []style.TextStyle{
	style.FgBlue,
	style.FgCyan,
	style.FgGreen,
	style.FgYellow,
	style.FgMagenta,
	style.FgRed,
}

func GetGraphiteStackListDisplayStrings(
	stacks []*models.GraphiteStack,
	getItemOperation func(item types.HasUrn) types.ItemOperation,
) [][]string {
	var displayStrings [][]string

	// Convert to the required format
	if len(stacks) > 0 {
		stack := stacks[len(stacks)-1]
		displayStrings = append(displayStrings, getGraphiteStackDisplayStrings(stack, getItemOperation(stack), 0)...)
	}

	return displayStrings
}

func GraphiteStackStatus(
	stack *models.GraphiteStack,
	itemOperation types.ItemOperation,
) string {
	if itemOperation != types.ItemOperationNone {
		return style.FgCyan.Sprintf("%s", ItemOperationToString(itemOperation, nil))
	}

	if stack.NeedsRestack {
		return style.FgYellow.Sprintf("%s", icons.STACK_ICON)
	}

	return ""
}

func getGraphiteStackDisplayStrings(
	node *models.GraphiteStack,
	itemOperation types.ItemOperation,
	depth int,
) [][]string {
	var result [][]string

	// Process children first
	for i, child := range node.Children {
		// Depth of children is dependent on the index of the child
		result = append(result, getGraphiteStackDisplayStrings(child, itemOperation, depth+i)...)
	}

	// Then process the current node
	stackStatus := GraphiteStackStatus(node, itemOperation)
	indicator := stackIndicatorEmpty
	if node.IsCurrent {
		indicator = stackIndicatorCurrent
	}
	if node.IsTrunk {
		indicator = style.FgMagenta.Sprint(indicator)
	} else {
		indicator = stackColors[depth%len(stackColors)].Sprint(indicator)
	}

	var indent string
	for i := 0; i < depth; i++ {
		indent += stackColors[i%len(stackColors)].Sprint(stackVerticalLine)
	}
	connectors := strings.Repeat(style.FgDefault.Sprint(stackHorizontalLine), max(0, len(node.Children)-2))
	if len(node.Children) > 1 {
		connectors += stackEndLine
	}

	line := indent + indicator + connectors + " " + node.Name
	if stackStatus != "" {
		line += " " + stackStatus
	}
	result = append(result, []string{line})

	return result
}
