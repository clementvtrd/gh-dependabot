package cmd

import (
	"fmt"
	"io"
	"os"
)

// Run dispatches the given arguments to the appropriate subcommand.
func Run(args []string) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printUsage(os.Stdout)
		return nil
	}

	sub, rest := args[0], args[1:]
	switch sub {
	case "approve":
		return runApprove(rest, os.Stdout)
	case "merge":
		return runMerge(rest, os.Stdout)
	case "rebase":
		return runRebase(rest, os.Stdout)
	default:
		return fmt.Errorf("unknown command %q\n\nRun 'gh dependabot --help' for usage", sub)
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: gh dependabot <command> [flags]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  approve   Approve all passing Dependabot PRs")
	fmt.Fprintln(w, "  merge     Merge all passing Dependabot PRs")
	fmt.Fprintln(w, "  rebase    Ask Dependabot to rebase all its open PRs")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  --dry-run   Print actions without executing them")
	fmt.Fprintln(w, "  --rebase    (merge only) Use rebase merge strategy")
}
