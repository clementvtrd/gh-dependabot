package cmd

import (
	"flag"
	"fmt"
	"io"
)

func runRebase(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("rebase", flag.ContinueOnError)
	fs.SetOutput(out)
	dryRun := fs.Bool("dry-run", false, "Print actions without executing them")
	if err := fs.Parse(args); err != nil {
		return err
	}

	prs, err := listDependabotPRs(false)
	if err != nil {
		return err
	}

	if len(prs) == 0 {
		fmt.Fprintln(out, "No eligible Dependabot PRs found.")
		return nil
	}

	for _, number := range prs {
		if *dryRun {
			fmt.Fprintf(out, "[dry-run] Would rebase PR #%d\n", number)
			continue
		}
		fmt.Fprintf(out, "Asking Dependabot to rebase PR #%d\n", number)
		if err := ghRun("pr", "comment", fmt.Sprintf("%d", number), "--body", "@dependabot rebase"); err != nil {
			return fmt.Errorf("posting rebase comment on PR #%d: %w", number, err)
		}
	}
	return nil
}
