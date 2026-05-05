package cmd

import (
	"flag"
	"fmt"
	"io"
)

func runApprove(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("approve", flag.ContinueOnError)
	fs.SetOutput(out)
	dryRun := fs.Bool("dry-run", false, "Print actions without executing them")
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

	for _, number := range prs {
		if *dryRun {
			fmt.Fprintf(out, "[dry-run] Would approve PR #%d\n", number)
			continue
		}
		fmt.Fprintf(out, "Approving PR #%d\n", number)
		if err := ghRun("pr", "review", fmt.Sprintf("%d", number), "--approve"); err != nil {
			return fmt.Errorf("approving PR #%d: %w", number, err)
		}
	}
	return nil
}
