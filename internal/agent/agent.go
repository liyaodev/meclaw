// Package agent routes work to ACP / CLI / HTTP / OpenAI-compatible agents.
package agent

import (
	"context"
	"fmt"
	"os"

	"github.com/meclaw/meclaw/internal/config"
)

// Request is a normalized prompt to an agent backend.
type Request struct {
	AgentID  string
	Session  string
	Prompt   string
	CWD      string
	Messages []ChatMessage // optional; used by openai mode
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
	openai *OpenAIRunner
}

// NewRouter builds a Router with default runners.
func NewRouter(agents map[string]config.Agent) *Router {
	return &Router{
		agents: agents,
		cli:    NewCLIRunner(),
		http:   NewHTTPRunner(nil),
		acp:    NewACPRunner(),
		openai: NewOpenAIRunner(nil),
	}
}

// WithHTTPClient replaces the HTTP runner's client (tests).
func (r *Router) WithHTTPClient(c *HTTPRunner) *Router {
	if c != nil {
		r.http = c
	}
	return r
}

// WithOpenAIRunner replaces the OpenAI runner (tests).
func (r *Router) WithOpenAIRunner(o *OpenAIRunner) *Router {
	if o != nil {
		r.openai = o
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
	case "openai":
		apiKey := cfg.APIKey
		if apiKey == "" {
			apiKey = os.Getenv("MECLAW_OPENAI_API_KEY")
		}
		msgs := req.Messages
		if len(msgs) == 0 {
			msgs = []ChatMessage{{Role: "user", Content: req.Prompt}}
		}
		return r.openai.RunChat(ctx, cfg.BaseURL, apiKey, cfg.Model, msgs)
	default:
		return Response{}, fmt.Errorf("agent %q: unsupported mode %q", req.AgentID, cfg.Mode)
	}
}
