package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
)

func GetGraphiteStackListDisplayStrings(
	stacks []*models.GraphiteStack,
	getItemOperation func(item types.HasUrn) types.ItemOperation,
	fullDescription bool,
	viewWidth int,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
) [][]string {
	var displayStrings [][]string

	// Flat list of all branches, categorized by type
	var trunkBranches []*models.GraphiteStack
	var regularBranches []*models.GraphiteStack

	for _, stack := range stacks {
		collectBranches(stack, &trunkBranches, &regularBranches)
	}

	// Format in the style of 'gt log short'
	if len(regularBranches) > 0 {
		// First regular branch
		branch := regularBranches[0]
		displayStrings = append(displayStrings, []string{
			branchDisplay(branch),
		})

		// Second branch with vertical connection
		if len(regularBranches) > 1 {
			displayStrings = append(displayStrings, []string{
				style.FgDefault.Sprint("│ ") + branchDisplay(regularBranches[1]),
			})
		}
	}

	// Trunk branch at the bottom with special connector
	if len(trunkBranches) > 0 {
		displayStrings = append(displayStrings, []string{
			style.FgBlue.Sprintf("◯─┘  %s", trunkBranches[0].Name),
		})
	}

	return displayStrings
}

// collectBranches traverses the stack and collects all branches into appropriate slices
func collectBranches(stack *models.GraphiteStack, trunkBranches, regularBranches *[]*models.GraphiteStack) {
	if stack == nil {
		return
	}

	// Add this branch to the appropriate slice
	if stack.IsTrunk {
		*trunkBranches = append(*trunkBranches, stack)
	} else {
		*regularBranches = append(*regularBranches, stack)
	}

	// Recursively process children
	for _, child := range stack.Children {
		collectBranches(child, trunkBranches, regularBranches)
	}

	// Also process parents to ensure we get all branches
	for _, parent := range stack.Parents {
		// Only process parents that aren't already in our lists
		parentInTrunk := false
		for _, trunk := range *trunkBranches {
			if trunk.Name == parent.Name {
				parentInTrunk = true
				break
			}
		}

		parentInRegular := false
		for _, regular := range *regularBranches {
			if regular.Name == parent.Name {
				parentInRegular = true
				break
			}
		}

		if !parentInTrunk && !parentInRegular {
			collectBranches(parent, trunkBranches, regularBranches)
		}
	}
}

// branchDisplay formats a branch with appropriate styling
func branchDisplay(branch *models.GraphiteStack) string {
	if branch.IsCurrent {
		return style.FgYellow.Sprintf("◉    %s", branch.Name)
	}
	return style.FgGreen.Sprintf("◯    %s", branch.Name)
}
