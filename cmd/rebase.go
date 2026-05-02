package cmd

import (
	"fmt"

	gh "github.com/cli/go-gh/v2"
	"github.com/spf13/cobra"
)

func newRebaseCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "rebase",
		Short: "Ask Dependabot to rebase all its open PRs",
		Long: `Post "@dependabot rebase" on every open Dependabot pull request,
regardless of their CI check status.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			prs, err := listDependabotPRs(false)
			if err != nil {
				return err
			}

			if len(prs) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No eligible Dependabot PRs found.")
				return nil
			}

			for _, number := range prs {
				prStr := fmt.Sprintf("%d", number)
				if dryRun {
					fmt.Fprintf(cmd.OutOrStdout(), "[dry-run] Would rebase PR #%d\n", number)
					continue
				}

				fmt.Fprintf(cmd.OutOrStdout(), "Asking Dependabot to rebase PR #%d\n", number)
				if _, _, err := gh.Exec("pr", "comment", prStr, "--body", "@dependabot rebase"); err != nil {
					return fmt.Errorf("posting rebase comment on PR #%d: %w", number, err)
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print actions without executing them")
	return cmd
}
