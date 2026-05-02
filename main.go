package main

import (
	"os"

	"github.com/clementvtrd/gh-dependabot/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
