package cmd

import (
	"fmt"
	"os"

	"github.com/sjramblings/yaml-merge/internal/merger"
	"github.com/sjramblings/yaml-merge/internal/progress"
	"github.com/spf13/cobra"
)

var (
	// Version information
	version   string
	gitCommit string
	buildTime string

	// CLI flags
	file1 string
	file2 string
	key   string
)

var rootCmd = &cobra.Command{
	Use:   "yaml-merge",
	Short: "A tool to merge YAML files",
	Long: `yaml-merge is a CLI tool that merges two YAML files based on a specified key.
It can combine sequences from both files while removing duplicates.`,
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create progress tracker
		prog := &progress.ConsoleProgress{}

		// Validate input files exist
		if _, err := os.Stat(file1); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", file1)
		}
		if _, err := os.Stat(file2); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", file2)
		}

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
	// Required flags
	rootCmd.Flags().StringVarP(&file1, "file1", "1", "", "First YAML file to merge (required)")
	rootCmd.Flags().StringVarP(&file2, "file2", "2", "", "Second YAML file to merge (required)")
	rootCmd.Flags().StringVarP(&key, "key", "k", "", "Key to merge on (required)")

	// Mark flags as required
	rootCmd.MarkFlagRequired("file1")
	rootCmd.MarkFlagRequired("file2")
	rootCmd.MarkFlagRequired("key")
}
