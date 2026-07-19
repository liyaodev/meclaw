package agent

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// CLIRunner executes a local command and returns stdout as the reply.
type CLIRunner struct{}

// NewCLIRunner creates a CLIRunner.
func NewCLIRunner() *CLIRunner {
	return &CLIRunner{}
}

// Run is unused on CLIRunner directly; use RunCommand via Router.
func (r *CLIRunner) Run(ctx context.Context, req Request) (Response, error) {
	return Response{}, fmt.Errorf("cli runner requires command via RunCommand")
}

// RunCommand runs command with args, appending the prompt as a final argument.
func (r *CLIRunner) RunCommand(ctx context.Context, command string, args []string, req Request) (Response, error) {
	if command == "" {
		return Response{}, fmt.Errorf("cli: empty command")
	}
	fullArgs := append(append([]string{}, args...), req.Prompt)
	cmd := exec.CommandContext(ctx, command, fullArgs...)
	if req.CWD != "" {
		cmd.Dir = req.CWD
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return Response{}, fmt.Errorf("cli %s: %s", command, errMsg)
	}
	return Response{Text: strings.TrimSpace(stdout.String())}, nil
}
