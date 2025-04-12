package graphite_commands

import (
	"encoding/json"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

// Repo contains methods for interacting with Graphite repositories
type Repo struct {
	*Common
}

// NewRepo creates a new Repo instance
func NewRepo(common *Common) *Repo {
	return &Repo{
		Common: common,
	}
}

// Init initializes Graphite in a repository
func (r *Repo) Init(trunkName string) error {
	cmdArgs := NewGraphiteCmd("init").
		ArgIf(trunkName != "", "--trunk", trunkName).
		ToArgv()

	return r.GetCmd().New(cmdArgs).Run()
}

// Sync pulls in new changes from the remote repository
func (r *Repo) Sync() error {
	cmdArgs := NewGraphiteCmd("sync").
		ToArgv()

	return r.GetCmd().New(cmdArgs).Run()
}

// Status shows the current status of the repository
func (r *Repo) Status() error {
	cmdArgs := NewGraphiteCmd("status").
		ToArgv()

	return r.GetCmd().New(cmdArgs).Run()
}

// Auth authenticates with Graphite
func (r *Repo) Auth() error {
	cmdArgs := NewGraphiteCmd("auth").
		ToArgv()

	return r.GetCmd().New(cmdArgs).Run()
}

// User returns information about the current user
func (r *Repo) User() error {
	cmdArgs := NewGraphiteCmd("user").
		ToArgv()

	return r.GetCmd().New(cmdArgs).Run()
}

// State returns the current state of all stacks
func (r *Repo) State() ([]*models.GraphiteStack, error) {
	cmdArgs := NewGraphiteCmd("state").
		ToArgv()

	output, err := r.GetCmd().New(cmdArgs).RunWithOutput()
	if err != nil {
		return nil, err
	}

	// Get current branch name using git_commands
	branchCommands := git_commands.NewBranchCommands(r.GitCommon)
	currentBranch, err := branchCommands.CurrentBranchName()
	if err != nil {
		return nil, err
	}

	// Parse the JSON output
	var state map[string]struct {
		Trunk        bool `json:"trunk"`
		NeedsRestack bool `json:"needs_restack"`
		Parents      []struct {
			Ref string `json:"ref"`
			SHA string `json:"sha"`
		} `json:"parents"`
	}

	if err := json.Unmarshal([]byte(output), &state); err != nil {
		return nil, err
	}

	// First pass: create all stacks
	stacksByName := make(map[string]*models.GraphiteStack)
	for branchName, branchInfo := range state {
		stack := &models.GraphiteStack{
			Name:      branchName,
			IsCurrent: branchName == currentBranch,
			IsTrunk:   branchInfo.Trunk,
			Children:  make([]*models.GraphiteStack, 0),
		}
		stacksByName[branchName] = stack
	}

	// Second pass: establish parent-child relationships
	for branchName, branchInfo := range state {
		stack := stacksByName[branchName]

		// Add parent references
		for _, parent := range branchInfo.Parents {
			if parentStack, exists := stacksByName[parent.Ref]; exists {
				stack.Parents = append(stack.Parents, parentStack)
				parentStack.Children = append(parentStack.Children, stack)
			}
		}
	}

	// Convert map to slice, starting with trunk branches
	var stacks []*models.GraphiteStack
	for _, stack := range stacksByName {
		// Only include stacks that have no parents (trunk branches) or are not children of any other stack
		isChild := false
		for _, otherStack := range stacksByName {
			for _, child := range otherStack.Children {
				if child == stack {
					isChild = true
					break
				}
			}
			if isChild {
				break
			}
		}

		if !isChild || stack.IsTrunk {
			stacks = append(stacks, stack)
		}
	}

	return stacks, nil
}
