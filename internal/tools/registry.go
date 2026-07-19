// Package tools registers and invokes agent tools (scenario A3).
package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/meclaw/meclaw/internal/policy"
	"github.com/meclaw/meclaw/internal/sandbox"
)

// Spec describes a callable tool.
type Spec struct {
	Name        string
	Description string
}

// CallRequest is an invocation request.
type CallRequest struct {
	Name   string
	UserID string
	Args   map[string]any
}

// CallResult is a tool result.
type CallResult struct {
	Output string
}

// Registry lists and runs tools.
type Registry interface {
	List() []Spec
	Call(ctx context.Context, req CallRequest) (CallResult, error)
}

// Handler implements one tool.
type Handler func(ctx context.Context, args map[string]any) (string, error)

// MapRegistry is an in-memory tool registry with policy checks.
type MapRegistry struct {
	specs    []Spec
	handlers map[string]Handler
	policy   policy.Engine
}

// NewMapRegistry creates a registry. pol may be nil (allow all tools).
func NewMapRegistry(pol policy.Engine) *MapRegistry {
	if pol == nil {
		pol = policy.NewAllowList(nil, nil)
	}
	return &MapRegistry{
		handlers: map[string]Handler{},
		policy:   pol,
	}
}

// Register adds a tool.
func (r *MapRegistry) Register(spec Spec, h Handler) {
	r.specs = append(r.specs, spec)
	r.handlers[spec.Name] = h
}

// List returns registered tools.
func (r *MapRegistry) List() []Spec {
	out := make([]Spec, len(r.specs))
	copy(out, r.specs)
	return out
}

// Call invokes a tool after policy check.
func (r *MapRegistry) Call(ctx context.Context, req CallRequest) (CallResult, error) {
	if !r.policy.AllowTool(req.UserID, req.Name) {
		return CallResult{}, fmt.Errorf("tool %q denied by policy", req.Name)
	}
	h, ok := r.handlers[req.Name]
	if !ok {
		return CallResult{}, fmt.Errorf("unknown tool %q", req.Name)
	}
	out, err := h(ctx, req.Args)
	if err != nil {
		return CallResult{}, err
	}
	return CallResult{Output: out}, nil
}

// DefaultRegistry registers built-in echo + shell (via sandbox).
func DefaultRegistry(pol policy.Engine, exec sandbox.Executor) *MapRegistry {
	r := NewMapRegistry(pol)
	r.Register(Spec{Name: "echo", Description: "echo text argument"}, func(ctx context.Context, args map[string]any) (string, error) {
		_ = ctx
		return fmt.Sprint(args["text"]), nil
	})
	r.Register(Spec{Name: "shell", Description: "run whitelisted command via sandbox"}, func(ctx context.Context, args map[string]any) (string, error) {
		cmd, _ := args["command"].(string)
		if cmd == "" {
			return "", fmt.Errorf("shell: command required")
		}
		var argv []string
		switch v := args["args"].(type) {
		case []string:
			argv = v
		case []any:
			for _, a := range v {
				argv = append(argv, fmt.Sprint(a))
			}
		case string:
			if v != "" {
				argv = strings.Fields(v)
			}
		}
		res, err := exec.Exec(ctx, sandbox.ExecRequest{Command: cmd, Args: argv})
		if err != nil {
			return "", err
		}
		if res.Stderr != "" {
			return res.Stdout + "\n" + res.Stderr, nil
		}
		return res.Stdout, nil
	})
	return r
}
