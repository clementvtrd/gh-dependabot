package cmd

import (
	"fmt"

	gh "github.com/cli/go-gh/v2"
	"github.com/spf13/cobra"
)

func newApproveCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Approve all passing Dependabot PRs",
		Long: `Approve every open Dependabot pull request whose CI checks have
all succeeded or were skipped (or that have no checks at all).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			prs, err := listDependabotPRs(true)
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
					fmt.Fprintf(cmd.OutOrStdout(), "[dry-run] Would approve PR #%d\n", number)
					continue
				}

				fmt.Fprintf(cmd.OutOrStdout(), "Approving PR #%d\n", number)
				if _, _, err := gh.Exec("pr", "review", prStr, "--approve"); err != nil {
					return fmt.Errorf("approving PR #%d: %w", number, err)
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print actions without executing them")
	return cmd
}
