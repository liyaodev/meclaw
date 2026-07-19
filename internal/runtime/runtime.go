// Package runtime wires policy, session, memory, tools, and agent routing (A1–A5).
package runtime

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/meclaw/meclaw/internal/agent"
	"github.com/meclaw/meclaw/internal/config"
	"github.com/meclaw/meclaw/internal/gateway"
	"github.com/meclaw/meclaw/internal/memory"
	"github.com/meclaw/meclaw/internal/observe"
	"github.com/meclaw/meclaw/internal/orchestrate"
	"github.com/meclaw/meclaw/internal/policy"
	"github.com/meclaw/meclaw/internal/sandbox"
	"github.com/meclaw/meclaw/internal/session"
	"github.com/meclaw/meclaw/internal/tools"
)

// Runtime handles inbound IM messages end-to-end.
type Runtime struct {
	cfg      *config.Config
	policy   policy.Engine
	audit    policy.Auditor
	store    session.StateStore
	mem      memory.Store
	router   *agent.Router
	tools    tools.Registry
	resolver orchestrate.Resolver
	skills   orchestrate.SkillLoader
	tracer   observe.Tracer
	idGen    func() string
}

// Options configures Runtime construction.
type Options struct {
	Policy   policy.Engine
	Audit    policy.Auditor
	Store    session.StateStore
	Memory   memory.Store
	Router   *agent.Router
	Tools    tools.Registry
	Resolver orchestrate.Resolver
	Tracer   observe.Tracer
	IDGen    func() string
}

// New builds a Runtime from config and optional overrides.
func New(cfg *config.Config, opts Options) *Runtime {
	store := opts.Store
	if store == nil {
		store = session.NewMemoryStore()
	}
	pol := opts.Policy
	if pol == nil {
		pol = policy.NewAllowList(cfg.Policy.AllowUsers, cfg.Policy.AllowTools)
	}
	audit := opts.Audit
	if audit == nil {
		audit = policy.NewWriterAuditor(os.Stdout)
	}
	router := opts.Router
	if router == nil {
		router = agent.NewRouter(cfg.Agents)
	}
	idGen := opts.IDGen
	if idGen == nil {
		idGen = newSessionID
	}
	mem := opts.Memory
	if mem == nil {
		mem = memory.NopStore{}
	}
	toolReg := opts.Tools
	if toolReg == nil {
		ex := sandbox.NewLocalExecutor(cfg.Sandbox.AllowCommands)
		toolReg = tools.DefaultRegistry(pol, ex)
	}
	resolver := opts.Resolver
	if resolver == nil {
		resolver = orchestrate.NewRuleResolver(cfg.Bindings)
	}
	tracer := opts.Tracer
	if tracer == nil {
		tracer = observe.NopTracer{}
	}
	return &Runtime{
		cfg:      cfg,
		policy:   pol,
		audit:    audit,
		store:    store,
		mem:      mem,
		router:   router,
		tools:    toolReg,
		resolver: resolver,
		skills:   orchestrate.SkillLoader{Root: cfg.SkillsDir},
		tracer:   tracer,
		idGen:    idGen,
	}
}

// NewFromConfig builds Runtime with file session, memory, and traces under cfg.DataDir.
func NewFromConfig(cfg *config.Config, opts Options) (*Runtime, error) {
	if opts.Store == nil {
		sessDir := filepath.Join(cfg.DataDir, "sessions")
		store, err := session.NewFileStore(sessDir)
		if err != nil {
			return nil, err
		}
		opts.Store = store
	}
	if opts.Memory == nil {
		if cfg.Memory.Enabled {
			memDir := filepath.Join(cfg.DataDir, "memory")
			mem, err := memory.NewFileStore(memDir)
			if err != nil {
				return nil, err
			}
			opts.Memory = mem
		} else {
			opts.Memory = memory.NopStore{}
		}
	}
	if opts.Tracer == nil {
		tr, err := observe.NewJSONLTracer(filepath.Join(cfg.DataDir, "traces"))
		if err != nil {
			return nil, err
		}
		opts.Tracer = tr
	}
	return New(cfg, opts), nil
}

func newSessionID() string {
	var b [16]byte
	if _, err := io.ReadFull(rand.Reader, b[:]); err != nil {
		return fmt.Sprintf("sess-%d", len(b))
	}
	return hex.EncodeToString(b[:])
}

// Handle processes one inbound message and returns a reply.
func (rt *Runtime) Handle(ctx context.Context, msg gateway.Message) (string, error) {
	span := rt.tracer.Start("runtime.handle", map[string]string{
		"channel": msg.Channel,
		"user":    msg.UserID,
		"chat":    msg.ChatID,
	})
	reply, err := rt.handle(ctx, msg)
	rt.tracer.End(span, err)
	return reply, err
}

func (rt *Runtime) handle(ctx context.Context, msg gateway.Message) (string, error) {
	key := sessionKey(msg)

	if !rt.policy.AllowMessage(msg) {
		rt.audit.Log(policy.Event{
			Action:  "deny",
			Channel: msg.Channel,
			UserID:  msg.UserID,
			ChatID:  msg.ChatID,
			Detail:  "user not allowed",
		})
		return "", fmt.Errorf("message denied by policy")
	}

	prompt := strings.TrimSpace(msg.Text)

	// A3: /tool <name> key=value ...
	if strings.HasPrefix(prompt, "/tool ") {
		return rt.handleTool(ctx, msg, prompt)
	}

	// A5: /skill <path> — attach skill text as one-shot system prefix via echo-back
	if strings.HasPrefix(prompt, "/skill ") {
		return rt.handleSkill(ctx, msg, prompt)
	}

	agentID := rt.resolveAgent(msg, key)

	if strings.HasPrefix(prompt, "/agent ") {
		parts := strings.Fields(prompt)
		if len(parts) >= 2 {
			want := parts[1]
			if _, ok := rt.cfg.Agents[want]; !ok {
				return "", fmt.Errorf("unknown agent %q", want)
			}
			rt.store.SetAgent(key, want)
			agentID = want
			prompt = strings.TrimSpace(strings.TrimPrefix(prompt, parts[0]+" "+parts[1]))
			if prompt == "" {
				rt.audit.Log(policy.Event{
					Action:  "allow",
					Channel: msg.Channel,
					UserID:  msg.UserID,
					ChatID:  msg.ChatID,
					Agent:   agentID,
					Detail:  "agent switched",
				})
				return fmt.Sprintf("switched to agent %q", agentID), nil
			}
		}
	}

	sessionID, ok := rt.store.Get(key)
	if !ok {
		sessionID = rt.idGen()
		rt.store.Set(key, sessionID)
	}

	history, err := rt.mem.Recall(ctx, key, rt.cfg.Memory.MaxMessages)
	if err != nil {
		return "", fmt.Errorf("memory recall: %w", err)
	}

	skillPrefix := ""
	if acfg, ok := rt.cfg.Agents[agentID]; ok && acfg.Skill != "" {
		text, loadErr := rt.skills.Load(acfg.Skill)
		if loadErr == nil {
			skillPrefix = text
		}
	}

	req := agent.Request{
		AgentID: agentID,
		Session: sessionID,
		Prompt:  prompt,
	}
	if cfg, ok := rt.cfg.Agents[agentID]; ok && cfg.Mode == "openai" {
		msgs := buildChatMessages(history, prompt)
		if skillPrefix != "" {
			msgs = append([]agent.ChatMessage{{Role: "system", Content: skillPrefix}}, msgs...)
		}
		req.Messages = msgs
	} else {
		if skillPrefix != "" {
			prompt = "Skill instructions:\n" + skillPrefix + "\n\nUser: " + prompt
		}
		if len(history) > 0 {
			req.Prompt = injectMemoryPrompt(history, prompt)
		} else {
			req.Prompt = prompt
		}
	}

	rt.audit.Log(policy.Event{
		Action:  "allow",
		Channel: msg.Channel,
		UserID:  msg.UserID,
		ChatID:  msg.ChatID,
		Agent:   agentID,
		Detail:  "dispatch",
	})

	resp, err := rt.router.Run(ctx, req)
	if err != nil {
		rt.audit.Log(policy.Event{
			Action:  "error",
			Channel: msg.Channel,
			UserID:  msg.UserID,
			ChatID:  msg.ChatID,
			Agent:   agentID,
			Detail:  err.Error(),
		})
		return "", err
	}

	_ = rt.mem.Append(ctx, key, memory.Message{Role: "user", Content: strings.TrimSpace(msg.Text)})
	_ = rt.mem.Append(ctx, key, memory.Message{Role: "assistant", Content: resp.Text})

	rt.audit.Log(policy.Event{
		Action:  "reply",
		Channel: msg.Channel,
		UserID:  msg.UserID,
		ChatID:  msg.ChatID,
		Agent:   agentID,
		Detail:  truncate(resp.Text, 120),
	})
	return resp.Text, nil
}

func (rt *Runtime) handleTool(ctx context.Context, msg gateway.Message, prompt string) (string, error) {
	parts := strings.Fields(prompt)
	if len(parts) < 2 {
		return "", fmt.Errorf("usage: /tool <name> [key=value...]")
	}
	name := parts[1]
	args := map[string]any{}
	for _, p := range parts[2:] {
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			args["text"] = strings.Join(parts[2:], " ")
			break
		}
		args[k] = v
	}
	if _, has := args["text"]; !has && name == "echo" && len(parts) > 2 {
		args["text"] = strings.Join(parts[2:], " ")
	}
	// shell: /tool shell command=echo args=hello
	res, err := rt.tools.Call(ctx, tools.CallRequest{Name: name, UserID: msg.UserID, Args: args})
	if err != nil {
		rt.audit.Log(policy.Event{
			Action:  "deny",
			Channel: msg.Channel,
			UserID:  msg.UserID,
			ChatID:  msg.ChatID,
			Detail:  err.Error(),
		})
		return "", err
	}
	rt.audit.Log(policy.Event{
		Action:  "reply",
		Channel: msg.Channel,
		UserID:  msg.UserID,
		ChatID:  msg.ChatID,
		Detail:  "tool:" + name,
	})
	return res.Output, nil
}

func (rt *Runtime) handleSkill(ctx context.Context, msg gateway.Message, prompt string) (string, error) {
	_ = ctx
	parts := strings.Fields(prompt)
	if len(parts) < 2 {
		return "", fmt.Errorf("usage: /skill <path-under-skills_dir>")
	}
	rel := parts[1]
	text, err := rt.skills.Load(rel)
	if err != nil {
		return "", err
	}
	rest := strings.TrimSpace(strings.TrimPrefix(prompt, parts[0]+" "+parts[1]))
	if rest == "" {
		return fmt.Sprintf("skill %q loaded (%d bytes). Send a message or: /skill %s <prompt>", rel, len(text), rel), nil
	}
	// One-shot: inject skill then dispatch with same agent
	key := sessionKey(msg)
	agentID := rt.resolveAgent(msg, key)
	sessionID, ok := rt.store.Get(key)
	if !ok {
		sessionID = rt.idGen()
		rt.store.Set(key, sessionID)
	}
	full := "Skill instructions:\n" + text + "\n\nUser: " + rest
	resp, err := rt.router.Run(ctx, agent.Request{AgentID: agentID, Session: sessionID, Prompt: full})
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}

func (rt *Runtime) resolveAgent(msg gateway.Message, key string) string {
	if id, ok := rt.store.GetAgent(key); ok && id != "" {
		return id
	}
	if id, ok := rt.resolver.Resolve(msg.Channel, msg.ChatID, msg.UserID); ok {
		return id
	}
	return rt.cfg.DefaultAgent
}

func buildChatMessages(history []memory.Message, prompt string) []agent.ChatMessage {
	out := make([]agent.ChatMessage, 0, len(history)+1)
	for _, m := range history {
		out = append(out, agent.ChatMessage{Role: m.Role, Content: m.Content})
	}
	out = append(out, agent.ChatMessage{Role: "user", Content: prompt})
	return out
}

func injectMemoryPrompt(history []memory.Message, prompt string) string {
	var b strings.Builder
	b.WriteString("Prior context:\n")
	for _, m := range history {
		b.WriteString(m.Role)
		b.WriteString(": ")
		b.WriteString(m.Content)
		b.WriteByte('\n')
	}
	b.WriteString("\nUser: ")
	b.WriteString(prompt)
	return b.String()
}

func sessionKey(msg gateway.Message) string {
	return msg.Channel + ":" + msg.ChatID + ":" + msg.UserID
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
