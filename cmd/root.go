package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command for gh-dependabot.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "gh-dependabot",
		Short: "Interact with pull requests opened by Dependabot",
		Long: `gh-dependabot is a GitHub CLI extension for managing pull requests
opened by Dependabot. It supports approving, merging, and rebasing
Dependabot PRs in bulk.`,
	}

	root.AddCommand(newApproveCmd())
	root.AddCommand(newMergeCmd())
	root.AddCommand(newRebaseCmd())

	return root
}
