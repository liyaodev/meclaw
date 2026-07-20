// Package sandbox isolates tool and command execution (scenario A3).
//
// Local process sandbox is next; browser / computer-use remain later interfaces.
package sandbox

import "context"

// ExecRequest is a sandboxed command.
type ExecRequest struct {
	Command string
	Args    []string
	CWD     string
}

// ExecResult captures stdout/stderr.
type ExecResult struct {
	Stdout string
	Stderr string
	Code   int
}

// Executor runs commands in an isolation boundary.
type Executor interface {
	Exec(ctx context.Context, req ExecRequest) (ExecResult, error)
}

// BrowserController is a later interface for browser automation.
type BrowserController interface {
	Navigate(ctx context.Context, url string) error
}

// ComputerController is a later interface for OS-level control.
type ComputerController interface {
	Screenshot(ctx context.Context) ([]byte, error)
}

// StubExecutor rejects all commands until A3 lands.
type StubExecutor struct{}

// Exec returns not implemented.
func (StubExecutor) Exec(ctx context.Context, req ExecRequest) (ExecResult, error) {
	_ = ctx
	_ = req
	return ExecResult{}, errNotImplemented()
}
