package presentation

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
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
	fullDescription bool,
	viewWidth int,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
) [][]string {
	var displayStrings [][]string

	// Convert to the required format
	if len(stacks) > 0 {
		displayStrings = append(displayStrings, getGraphiteStackDisplayStrings(stacks[len(stacks)-1], 0)...)
	}

	return displayStrings
}

func getGraphiteStackDisplayStrings(node *models.GraphiteStack, depth int) [][]string {
	var result [][]string

	// Process children first
	for i, child := range node.Children {
		// Depth of children is dependent on the index of the child
		result = append(result, getGraphiteStackDisplayStrings(child, depth+i)...)
	}

	// Then process the current node
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
	result = append(result, []string{line})

	return result
}
