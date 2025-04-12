package graphite_commands

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

// Repo contains methods for interacting with Graphite repositories
type Repo struct {
	*Common
	Branch *git_commands.BranchCommands
}

// NewRepo creates a new Repo instance
func NewRepo(common *Common) *Repo {
	return &Repo{
		Common: common,
		Branch: git_commands.NewBranchCommands(common.GitCommon),
	}
}

// Init initializes Graphite in a repository
func (r *Repo) Init(trunkName string) error {
	cmdArgs := NewGraphiteCmd("init")
	if trunkName != "" {
		cmdArgs.Arg("--trunk", trunkName)
	}

	return r.GetCmd().New(cmdArgs.ToArgv()).Run()
}

// Sync pulls in new changes from the remote repository
func (r *Repo) Sync() error {
	cmdArgs := NewGraphiteCmd("sync")
	return r.GetCmd().New(cmdArgs.ToArgv()).Run()
}

// Status shows the current status of the repository
func (r *Repo) Status() error {
	cmdArgs := NewGraphiteCmd("status")
	return r.GetCmd().New(cmdArgs.ToArgv()).Run()
}

// Auth authenticates with Graphite
func (r *Repo) Auth() error {
	cmdArgs := NewGraphiteCmd("auth")
	return r.GetCmd().New(cmdArgs.ToArgv()).Run()
}

// User returns information about the current user
func (r *Repo) User() error {
	cmdArgs := NewGraphiteCmd("user")
	return r.GetCmd().New(cmdArgs.ToArgv()).Run()
}

// State returns a list of branches in scope of the current trunk in post-order traversal
func (r *Repo) State() ([]*models.GraphiteStack, error) {
	cmdArgs := NewGraphiteCmd("state")
	output, err := r.GetCmd().New(cmdArgs.ToArgv()).DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}

	// The actual structure of gt state output
	type ParentInfo struct {
		Ref string `json:"ref"`
		Sha string `json:"sha"`
	}

	type BranchInfo struct {
		Trunk        bool         `json:"trunk"`
		NeedsRestack bool         `json:"needs_restack"`
		Parents      []ParentInfo `json:"parents"`
	}

	var state map[string]BranchInfo
	if err := json.Unmarshal([]byte(output), &state); err != nil {
		return nil, err
	}

	// Get current branch using the BranchCommands implementation
	currentBranch, err := r.Branch.CurrentBranchName()
	if err != nil {
		return nil, err
	}
	currentBranch = strings.TrimSpace(currentBranch)

	// First pass: create all branch nodes
	branchMap := make(map[string]*models.GraphiteStack)
	for branchName, info := range state {
		if _, exists := branchMap[branchName]; !exists {
			branchMap[branchName] = &models.GraphiteStack{
				Name:      branchName,
				IsCurrent: branchName == currentBranch,
				IsTrunk:   info.Trunk,
				Children:  []*models.GraphiteStack{},
				Parents:   []*models.GraphiteStack{},
			}
		} else {
			branchMap[branchName].IsTrunk = info.Trunk
			branchMap[branchName].IsCurrent = branchName == currentBranch
		}
	}

	// Second pass: establish parent-child relationships
	for branchName, info := range state {
		branch := branchMap[branchName]

		// Process all parents from the array
		for _, parentInfo := range info.Parents {
			parentRef := parentInfo.Ref

			parent, exists := branchMap[parentRef]
			if !exists {
				// Create parent if it doesn't exist yet
				parent = &models.GraphiteStack{
					Name:     parentRef,
					Children: []*models.GraphiteStack{},
					Parents:  []*models.GraphiteStack{},
				}
				branchMap[parentRef] = parent
			}

			// Add this branch as a child of its parent
			parent.Children = append(parent.Children, branch)
			// Sort children by name
			sort.Slice(parent.Children, func(i, j int) bool {
				return parent.Children[i].Name < parent.Children[j].Name
			})

			// Add parent to this branch's parents
			branch.Parents = append(branch.Parents, parent)
		}
	}

	// Find the trunk branch for the current branch
	var relevantTrunk *models.GraphiteStack

	// Start from current branch and navigate up to find its trunk
	if current, exists := branchMap[currentBranch]; exists {
		// Find the trunk by navigating upward
		branch := current
		visited := make(map[string]bool)

		for {
			// Avoid infinite loops in case of circular references
			if visited[branch.Name] {
				break
			}
			visited[branch.Name] = true

			// If this is a trunk, we've found our relevant trunk
			if branch.IsTrunk {
				relevantTrunk = branch
				break
			}

			// If no more parents, break the loop
			if len(branch.Parents) == 0 {
				break
			}

			// Move to the parent and continue looking
			branch = branch.Parents[0]
		}
	}

	// If we found a relevant trunk, return the post-order traversal from the trunk
	if relevantTrunk != nil {
		return flattenStackPostOrder(relevantTrunk), nil
	}

	return []*models.GraphiteStack{}, nil
}

func flattenStackPostOrder(stack *models.GraphiteStack) []*models.GraphiteStack {
	var result []*models.GraphiteStack
	for _, child := range stack.Children {
		result = append(result, flattenStackPostOrder(child)...)
	}

	result = append(result, stack)
	return result
}
