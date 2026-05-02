package cmd

import (
	"fmt"

	gh "github.com/cli/go-gh/v2"
	"github.com/spf13/cobra"
)

func newMergeCmd() *cobra.Command {
	var dryRun bool
	var useRebase bool

	cmd := &cobra.Command{
		Use:   "merge",
		Short: "Merge all passing Dependabot PRs",
		Long: `Merge every open Dependabot pull request whose CI checks have all
succeeded or were skipped (or that have no checks at all).

If a merge fails, gh-dependabot automatically posts "@dependabot rebase"
to ask Dependabot to rebase the PR before a future merge attempt.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			prs, err := listDependabotPRs(true)
			if err != nil {
				return err
			}

			if len(prs) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No eligible Dependabot PRs found.")
				return nil
			}

			mergeFlag := "-dm"
			if useRebase {
				mergeFlag = "-dr"
			}

			for _, number := range prs {
				prStr := fmt.Sprintf("%d", number)
				if dryRun {
					fmt.Fprintf(cmd.OutOrStdout(), "[dry-run] Would merge PR #%d\n", number)
					continue
				}

				fmt.Fprintf(cmd.OutOrStdout(), "Merging PR #%d\n", number)
				if _, _, err := gh.Exec("pr", "merge", prStr, mergeFlag); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Merge failed for PR #%d -> asking Dependabot to rebase\n", number)
					if _, _, commentErr := gh.Exec("pr", "comment", prStr, "--body", "@dependabot rebase"); commentErr != nil {
						return fmt.Errorf("posting rebase comment on PR #%d: %w", number, commentErr)
					}
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print actions without executing them")
	cmd.Flags().BoolVar(&useRebase, "rebase", false, "Use rebase merge strategy instead of merge commit")
	return cmd
}
