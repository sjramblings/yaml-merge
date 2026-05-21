package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootCommandQuietFlagSuppressesProgressOnStdout(t *testing.T) {
	tmpDir := t.TempDir()
	file1 := filepath.Join(tmpDir, "file1.yaml")
	file2 := filepath.Join(tmpDir, "file2.yaml")

	if err := os.WriteFile(file1, []byte(`
workloadAccounts:
  - name: Network-Dev
    email: dev@example.com
`), 0644); err != nil {
		t.Fatalf("failed to write file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte(`
workloadAccounts:
  - name: Network-Prod
    email: prod@example.com
`), 0644); err != nil {
		t.Fatalf("failed to write file2: %v", err)
	}

	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	os.Stdout = writer
	defer func() {
		os.Stdout = oldStdout
		rootCmd.SetArgs(nil)
		quiet = false
	}()

	rootCmd.SetArgs([]string{file1, file2, "workloadAccounts", "-q"})
	execErr := rootCmd.Execute()

	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close stdout writer: %v", err)
	}
	os.Stdout = oldStdout
	if execErr != nil {
		t.Fatalf("rootCmd.Execute() returned error: %v", execErr)
	}

	captured, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read captured stdout: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("failed to close stdout reader: %v", err)
	}

	output := string(captured)
	if !strings.Contains(output, "workloadAccounts:") {
		t.Fatalf("stdout = %q, want merged YAML containing workloadAccounts", output)
	}
	if strings.Contains(output, "=== Starting") {
		t.Fatalf("stdout should not contain progress output: %q", output)
	}
}
