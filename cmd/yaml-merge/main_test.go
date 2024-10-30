package main

import (
	"os"
	"testing"
)

var exitCode int
var exitFunc = func(code int) {
	exitCode = code
}

func TestRunWithNoArgs(t *testing.T) {
	oldStdout := os.Stdout
	w, _, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	origExit := exit
	exit = exitFunc
	defer func() { exit = origExit }()

	main()

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

func TestRunWithInvalidArgs(t *testing.T) {
	oldStdout := os.Stdout
	w, _, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	oldArgs := os.Args
	os.Args = []string{"cmd", "invalid"}
	defer func() { os.Args = oldArgs }()

	origExit := exit
	exit = exitFunc
	defer func() { exit = origExit }()

	main()

	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}
