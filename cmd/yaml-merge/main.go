package main

import (
	"fmt"
	"os"

	"github.com/sjramblings/yaml-merge/cmd"
)

// Build information (populated by linker flags)
var (
	version   = "dev"
	gitCommit = "none"
	buildTime = "unknown"
)

func main() {
	if err := cmd.Execute(version, gitCommit, buildTime); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
