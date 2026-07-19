// Package agent routes work to ACP / CLI / HTTP agents.
package agent

import (
	"context"
	"fmt"

	"github.com/meclaw/meclaw/internal/config"
)

// Request is a normalized prompt to an agent backend.
type Request struct {
	AgentID string
	Session string
	Prompt  string
	CWD     string
}

// Response is a normalized agent reply.
type Response struct {
	Text string
}

// Runner executes agent requests.
type Runner interface {
	Run(ctx context.Context, req Request) (Response, error)
}

// Router dispatches requests to the correct Runner based on agent config.
type Router struct {
	agents map[string]config.Agent
	cli    *CLIRunner
	http   *HTTPRunner
	acp    *ACPRunner
}

// NewRouter builds a Router with default CLI/HTTP/ACP runners.
func NewRouter(agents map[string]config.Agent) *Router {
	return &Router{
		agents: agents,
		cli:    NewCLIRunner(),
		http:   NewHTTPRunner(nil),
		acp:    NewACPRunner(),
	}
}

// WithHTTPClient replaces the HTTP runner's client (tests).
func (r *Router) WithHTTPClient(c *HTTPRunner) *Router {
	if c != nil {
		r.http = c
	}
	return r
}

// Lookup returns agent config.
func (r *Router) Lookup(agentID string) (config.Agent, bool) {
	cfg, ok := r.agents[agentID]
	return cfg, ok
}

// Run selects the backend for AgentID and executes the request.
func (r *Router) Run(ctx context.Context, req Request) (Response, error) {
	cfg, ok := r.agents[req.AgentID]
	if !ok {
		return Response{}, fmt.Errorf("unknown agent %q", req.AgentID)
	}
	switch cfg.Mode {
	case "cli":
		return r.cli.RunCommand(ctx, cfg.Command, cfg.Args, req)
	case "http":
		return r.http.RunURL(ctx, cfg.BaseURL, req)
	case "acp":
		return r.acp.Run(ctx, req)
	default:
		return Response{}, fmt.Errorf("agent %q: unsupported mode %q", req.AgentID, cfg.Mode)
	}
}
