package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// LocalExecutor runs whitelisted commands via exec.CommandContext.
type LocalExecutor struct {
	Allow map[string]struct{}
	Timeout time.Duration
}

// NewLocalExecutor builds an executor. Empty allow means deny-all (safe default for tools).
// Pass explicit allow list from config (e.g. echo, date, uname).
func NewLocalExecutor(allow []string) *LocalExecutor {
	m := make(map[string]struct{}, len(allow))
	for _, c := range allow {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		// Normalize to base name so /bin/echo and echo match.
		m[filepath.Base(c)] = struct{}{}
		m[c] = struct{}{}
	}
	return &LocalExecutor{Allow: m, Timeout: 30 * time.Second}
}

// Exec runs the command if allowed.
func (e *LocalExecutor) Exec(ctx context.Context, req ExecRequest) (ExecResult, error) {
	base := filepath.Base(req.Command)
	if _, ok := e.Allow[req.Command]; !ok {
		if _, ok := e.Allow[base]; !ok {
			return ExecResult{Code: -1}, fmt.Errorf("sandbox: command %q not allowed", req.Command)
		}
	}
	timeout := e.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, req.Command, req.Args...)
	if req.CWD != "" {
		cmd.Dir = req.CWD
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	res := ExecResult{
		Stdout: strings.TrimSpace(stdout.String()),
		Stderr: strings.TrimSpace(stderr.String()),
	}
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			res.Code = ee.ExitCode()
		} else {
			res.Code = -1
		}
		if res.Stderr == "" {
			res.Stderr = err.Error()
		}
		return res, fmt.Errorf("sandbox: %s: %s", req.Command, res.Stderr)
	}
	res.Code = 0
	return res, nil
}
