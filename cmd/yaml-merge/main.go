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

var exit = os.Exit

func main() {
	if err := cmd.Execute(version, gitCommit, buildTime); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		exit(1)
	}
}
