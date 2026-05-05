package main

import (
	"fmt"
	"os"

	"github.com/clementvtrd/gh-dependabot/cmd"
)

func main() {
	if err := cmd.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
