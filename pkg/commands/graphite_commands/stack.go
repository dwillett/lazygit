package graphite_commands

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

// Stack contains methods for interacting with Graphite stacks
type Stack struct {
	*Common
}

// NewStack creates a new Stack instance
func NewStack(common *Common) *Stack {
	return &Stack{
		Common: common,
	}
}

// List returns all Graphite stacks
func (s *Stack) List() ([]*models.GraphiteStack, error) {
	// Run the Graphite CLI command to list stacks
	cmdArgs := NewGraphiteCmd("log").
		ToArgv()

	output, err := s.GetCmd().New(cmdArgs).RunWithOutput()
	if err != nil {
		return nil, err
	}

	// Parse the output into GraphiteStack models
	var stacks []*models.GraphiteStack
	lines := strings.Split(output, "\n")

	// Track branches we've already added to avoid duplicates
	addedBranches := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip lines that don't represent branches (like timestamps or commit messages)
		if !strings.HasPrefix(line, "◉") && !strings.HasPrefix(line, "◯") && !strings.HasPrefix(line, "│") {
			continue
		}

		// Skip lines that are just separators or commit info
		if strings.HasPrefix(line, "│") && !strings.Contains(line, " - ") {
			continue
		}

		// Extract branch name and current status
		isCurrent := strings.HasPrefix(line, "◉")

		// Handle different line formats
		var branchName string
		if strings.HasPrefix(line, "◉") || strings.HasPrefix(line, "◯") {
			// Format: ◉ branch-name (current) or ◯ branch-name
			parts := strings.SplitN(line, " ", 2)
			if len(parts) < 2 {
				continue
			}

			// Remove the branch icon and "(current)" if present
			branchName = strings.TrimSuffix(parts[1], "(current)")
			branchName = strings.TrimSpace(branchName)
		} else if strings.Contains(line, " - ") {
			// Format: │ hash - commit message
			// This is a commit line, skip it
			continue
		} else {
			// Skip any other format
			continue
		}

		// Skip if we've already added this branch
		if addedBranches[branchName] {
			continue
		}

		// Add the branch to our tracking map
		addedBranches[branchName] = true

		stack := &models.GraphiteStack{
			Name:      branchName,
			IsCurrent: isCurrent,
		}

		stacks = append(stacks, stack)
	}

	return stacks, nil
}

// Submit submits a stack to the remote repository
func (s *Stack) Submit() error {
	cmdArgs := NewGraphiteCmd("submit").
		Arg("--stack").
		ToArgv()

	return s.GetCmd().New(cmdArgs).Run()
}

// Restack ensures all branches in a stack are based on the current version of their parents
func (s *Stack) Restack() error {
	cmdArgs := NewGraphiteCmd("restack").
		ToArgv()

	return s.GetCmd().New(cmdArgs).Run()
}

// Sync pulls in new changes from the main branch to open stacks
func (s *Stack) Sync() error {
	cmdArgs := NewGraphiteCmd("sync").
		ToArgv()

	return s.GetCmd().New(cmdArgs).Run()
}

// Test runs a command on each branch in the stack
func (s *Stack) Test(command string) error {
	cmdArgs := NewGraphiteCmd("stack").
		Arg("test", command).
		ToArgv()

	return s.GetCmd().New(cmdArgs).Run()
}

// Log shows the commit history for the current stack
func (s *Stack) Log(short bool) (string, error) {
	cmdArgs := NewGraphiteCmd("log").
		ArgIf(short, "short").
		ArgIf(!short, "long").
		ToArgv()

	return s.GetCmd().New(cmdArgs).RunWithOutput()
}
