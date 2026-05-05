package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// pr represents a GitHub pull request returned by `gh pr list`.
type pr struct {
	Number            int           `json:"number"`
	Author            prAuthor      `json:"author"`
	StatusCheckRollup []statusCheck `json:"statusCheckRollup"`
}

type prAuthor struct {
	Login string `json:"login"`
}

type statusCheck struct {
	Conclusion string `json:"conclusion"`
}

// ghRun executes a gh CLI command, wiring its stderr to os.Stderr.
func ghRun(args ...string) error {
	c := exec.Command("gh", args...)
	c.Stderr = os.Stderr
	return c.Run()
}

// ghOutput executes a gh CLI command and returns its stdout.
func ghOutput(args ...string) ([]byte, error) {
	return exec.Command("gh", args...).Output()
}

// isDependabotAuthor reports whether the given login belongs to Dependabot.
func isDependabotAuthor(login string) bool {
	return login == "app/dependabot" || login == "dependabot[bot]"
}

// allChecksPassedOrSkipped reports whether every check in the rollup has
// concluded with SUCCESS or SKIPPED (or there are no checks at all).
func allChecksPassedOrSkipped(checks []statusCheck) bool {
	for _, c := range checks {
		if c.Conclusion != "SUCCESS" && c.Conclusion != "SKIPPED" {
			return false
		}
	}
	return true
}

// listDependabotPRs returns PR numbers for open Dependabot pull requests.
// When requirePassingChecks is true, only PRs whose CI checks have all
// succeeded or were skipped are included.
func listDependabotPRs(requirePassingChecks bool) ([]int, error) {
	fields := "number,author"
	if requirePassingChecks {
		fields += ",statusCheckRollup"
	}

	out, err := ghOutput("pr", "list", "--json", fields)
	if err != nil {
		return nil, fmt.Errorf("listing pull requests: %w", err)
	}

	var prs []pr
	if err := json.Unmarshal(out, &prs); err != nil {
		return nil, fmt.Errorf("parsing pull request list: %w", err)
	}

	var numbers []int
	for _, p := range prs {
		if !isDependabotAuthor(p.Author.Login) {
			continue
		}
		if requirePassingChecks && !allChecksPassedOrSkipped(p.StatusCheckRollup) {
			continue
		}
		numbers = append(numbers, p.Number)
	}
	return numbers, nil
}
