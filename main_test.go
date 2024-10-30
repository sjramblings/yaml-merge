package main

import (
	"os"
	"testing"
)

func TestRunWithNoArgs(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	w, _, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	// Save original args and restore after test
	oldArgs := os.Args
	os.Args = []string{"yaml-merge"}
	defer func() { os.Args = oldArgs }()

	// Mock exit function
	exitCode := 0
	origExit := Exit
	Exit = func(code int) {
		exitCode = code
	}
	defer func() { Exit = origExit }()

	main()

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

func TestRunWithInvalidArgs(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	w, _, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	// Save original args and restore after test
	oldArgs := os.Args
	os.Args = []string{"yaml-merge", "invalid"}
	defer func() { os.Args = oldArgs }()

	// Mock exit function
	exitCode := 0
	origExit := Exit
	Exit = func(code int) {
		exitCode = code
	}
	defer func() { Exit = origExit }()

	main()

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}
