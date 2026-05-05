package cmd

import (
	"flag"
	"fmt"
	"io"
)

func runMerge(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	fs.SetOutput(out)
	dryRun := fs.Bool("dry-run", false, "Print actions without executing them")
	useRebase := fs.Bool("rebase", false, "Use rebase merge strategy instead of merge commit")
	if err := fs.Parse(args); err != nil {
		return err
	}

	prs, err := listDependabotPRs(true)
	if err != nil {
		return err
	}

	if len(prs) == 0 {
		fmt.Fprintln(out, "No eligible Dependabot PRs found.")
		return nil
	}

	mergeFlag := "-dm"
	if *useRebase {
		mergeFlag = "-dr"
	}

	for _, number := range prs {
		prStr := fmt.Sprintf("%d", number)
		if *dryRun {
			fmt.Fprintf(out, "[dry-run] Would merge PR #%d\n", number)
			continue
		}
		fmt.Fprintf(out, "Merging PR #%d\n", number)
		if err := ghRun("pr", "merge", prStr, mergeFlag); err != nil {
			fmt.Fprintf(out, "Merge failed for PR #%d -> asking Dependabot to rebase\n", number)
			if commentErr := ghRun("pr", "comment", prStr, "--body", "@dependabot rebase"); commentErr != nil {
				return fmt.Errorf("posting rebase comment on PR #%d: %w", number, commentErr)
			}
		}
	}
	return nil
}
