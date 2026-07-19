// Package runtime wires policy, session, and agent routing for scenario A.
package runtime

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/meclaw/meclaw/internal/agent"
	"github.com/meclaw/meclaw/internal/config"
	"github.com/meclaw/meclaw/internal/gateway"
	"github.com/meclaw/meclaw/internal/policy"
	"github.com/meclaw/meclaw/internal/session"
)

// Runtime handles inbound IM messages end-to-end.
type Runtime struct {
	cfg     *config.Config
	policy  policy.Engine
	audit   policy.Auditor
	store   *session.MemoryStore
	router  *agent.Router
	idGen   func() string
}

// Options configures Runtime construction.
type Options struct {
	Policy policy.Engine
	Audit  policy.Auditor
	Store  *session.MemoryStore
	Router *agent.Router
	IDGen  func() string
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
	return &Runtime{
		cfg:    cfg,
		policy: pol,
		audit:  audit,
		store:  store,
		router: router,
		idGen:  idGen,
	}
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
	agentID := rt.resolveAgent(key)

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

	rt.audit.Log(policy.Event{
		Action:  "allow",
		Channel: msg.Channel,
		UserID:  msg.UserID,
		ChatID:  msg.ChatID,
		Agent:   agentID,
		Detail:  "dispatch",
	})

	resp, err := rt.router.Run(ctx, agent.Request{
		AgentID: agentID,
		Session: sessionID,
		Prompt:  prompt,
	})
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

func (rt *Runtime) resolveAgent(key string) string {
	if id, ok := rt.store.GetAgent(key); ok && id != "" {
		return id
	}
	return rt.cfg.DefaultAgent
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
