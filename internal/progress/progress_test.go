package progress

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestNewConsoleWriterWritesToStderr(t *testing.T) {
	writer := NewConsoleWriter(false)

	if writer.writer != os.Stderr {
		t.Fatalf("NewConsoleWriter(false).writer = %v, want os.Stderr", writer.writer)
	}
}

func TestConsoleWriterQuietSuppressesOutput(t *testing.T) {
	var buf bytes.Buffer
	writer := &ConsoleWriter{writer: &buf, quiet: true}

	writer.Start("Starting YAML Merge")
	writer.Step("Processing files")
	writer.Success("Read file")
	writer.Error("Failed file")
	writer.End()

	if buf.String() != "" {
		t.Fatalf("quiet writer emitted output %q, want empty output", buf.String())
	}
}

func TestConsoleWriterNonQuietEmitsProgressOutput(t *testing.T) {
	var buf bytes.Buffer
	writer := &ConsoleWriter{writer: &buf, quiet: false}

	writer.Start("Starting YAML Merge")
	writer.Step("Processing %s", "file1.yaml")
	writer.Success("Read %s", "file1.yaml")
	writer.Error("Failed %s", "file2.yaml")
	writer.End()

	output := buf.String()
	if output == "" {
		t.Fatal("non-quiet writer emitted empty output")
	}
	for _, expected := range []string{
		"=== Starting YAML Merge ===",
		"Processing file1.yaml",
		"Read file1.yaml",
		"Failed file2.yaml",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got %q", expected, output)
		}
	}
}
