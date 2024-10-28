package progress

import (
	"fmt"
	"io"
	"os"
)

type Writer interface {
	Start(operation string)
	Step(format string, a ...interface{})
	Success(format string, a ...interface{})
	Error(format string, a ...interface{})
	End()
}

type ConsoleWriter struct {
	writer io.Writer
	quiet  bool
}

func NewConsoleWriter(quiet bool) *ConsoleWriter {
	return &ConsoleWriter{
		writer: os.Stdout,
		quiet:  quiet,
	}
}

func (w *ConsoleWriter) Start(operation string) {
	if !w.quiet {
		fmt.Fprintf(w.writer, "=== %s ===\n", operation)
	}
}

func (w *ConsoleWriter) Step(format string, a ...interface{}) {
	if !w.quiet {
		fmt.Fprintf(w.writer, "→ "+format+"\n", a...)
	}
}

func (w *ConsoleWriter) Success(format string, a ...interface{}) {
	if !w.quiet {
		fmt.Fprintf(w.writer, "✓ "+format+"\n", a...)
	}
}

func (w *ConsoleWriter) Error(format string, a ...interface{}) {
	if !w.quiet {
		fmt.Fprintf(w.writer, "✗ "+format+"\n", a...)
	}
}

func (w *ConsoleWriter) End() {
	if !w.quiet {
		fmt.Fprintln(w.writer)
	}
}
