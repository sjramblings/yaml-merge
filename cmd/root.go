package cmd

import (
	"fmt"

	"github.com/sjramblings/yaml-merge/internal/merger"
	"github.com/sjramblings/yaml-merge/internal/progress"
	"github.com/spf13/cobra"
)

var (
	// Version information
	version   string
	gitCommit string
	buildTime string
)

var rootCmd = &cobra.Command{
	Use:   "yaml-merge <file1> <file2> <key>",
	Short: "A tool to merge YAML files",
	Long: `yaml-merge is a CLI tool that merges two YAML files based on a specified key.
It can combine sequences from both files while removing duplicates.

Example:
  yaml-merge file1.yaml file2.yaml workloadAccounts`,
	Version: version,
	Args:    cobra.ExactArgs(3), // Require exactly 3 arguments
	RunE: func(cmd *cobra.Command, args []string) error {
		file1 := args[0]
		file2 := args[1]
		key := args[2]

		// Create progress tracker
		prog := progress.NewConsoleWriter(false)

		// Perform merge
		result, err := merger.MergeYAMLFiles(file1, file2, key, prog)
		if err != nil {
			return fmt.Errorf("merge failed: %w", err)
		}

		// Print result to stdout
		fmt.Println(string(result))
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(ver, commit, built string) error {
	// Set version information
	version = ver
	gitCommit = commit
	buildTime = built

	// Add version template
	rootCmd.SetVersionTemplate(`Version: {{.Version}}
Git Commit: ` + gitCommit + `
Build Time: ` + buildTime + "\n")

	return rootCmd.Execute()
}

func init() {
	// No flags needed anymore
}
