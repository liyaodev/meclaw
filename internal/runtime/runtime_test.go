package runtime_test

import (
	"context"
	"strings"
	"testing"

	"github.com/meclaw/meclaw/internal/agent"
	"github.com/meclaw/meclaw/internal/config"
	"github.com/meclaw/meclaw/internal/gateway"
	"github.com/meclaw/meclaw/internal/policy"
	"github.com/meclaw/meclaw/internal/runtime"
	"github.com/meclaw/meclaw/internal/session"
)

func TestHandleAllowAndSessionReuse(t *testing.T) {
	cfg := &config.Config{
		DefaultAgent: "echo",
		Agents: map[string]config.Agent{
			"echo": {Mode: "cli", Command: "echo"},
		},
	}
	audit := &policy.BufferAuditor{}
	store := session.NewMemoryStore()
	var ids []string
	rt := runtime.New(cfg, runtime.Options{
		Audit:  audit,
		Store:  store,
		Router: agent.NewRouter(cfg.Agents),
		IDGen: func() string {
			id := "sess-" + string(rune('a'+len(ids)))
			ids = append(ids, id)
			return id
		},
	})

	ctx := context.Background()
	msg := gateway.Message{Channel: "stdio", UserID: "u1", ChatID: "c1", Text: "hello"}
	reply, err := rt.Handle(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}
	if reply != "hello" {
		t.Fatalf("reply=%q", reply)
	}
	sid1, ok := store.Get("stdio:c1:u1")
	if !ok {
		t.Fatal("session missing")
	}

	_, err = rt.Handle(ctx, gateway.Message{Channel: "stdio", UserID: "u1", ChatID: "c1", Text: "again"})
	if err != nil {
		t.Fatal(err)
	}
	sid2, _ := store.Get("stdio:c1:u1")
	if sid1 != sid2 {
		t.Fatalf("session not reused: %s vs %s", sid1, sid2)
	}
	if len(ids) != 1 {
		t.Fatalf("id gen called %d times", len(ids))
	}
}

func TestHandleDeny(t *testing.T) {
	cfg := &config.Config{
		DefaultAgent: "echo",
		Agents:       map[string]config.Agent{"echo": {Mode: "cli", Command: "echo"}},
		Policy:       config.Policy{AllowUsers: []string{"alice"}},
	}
	audit := &policy.BufferAuditor{}
	rt := runtime.New(cfg, runtime.Options{Audit: audit})
	_, err := rt.Handle(context.Background(), gateway.Message{UserID: "bob", Text: "hi", Channel: "http", ChatID: "c"})
	if err == nil {
		t.Fatal("expected deny")
	}
	ev := audit.Snapshot()
	if len(ev) == 0 || ev[0].Action != "deny" {
		t.Fatalf("%+v", ev)
	}
}

func TestAgentSwitch(t *testing.T) {
	cfg := &config.Config{
		DefaultAgent: "echo",
		Agents: map[string]config.Agent{
			"echo":  {Mode: "cli", Command: "echo"},
			"echo2": {Mode: "cli", Command: "echo"},
		},
	}
	store := session.NewMemoryStore()
	rt := runtime.New(cfg, runtime.Options{Store: store, Audit: &policy.BufferAuditor{}})
	ctx := context.Background()
	msg := gateway.Message{Channel: "stdio", UserID: "u", ChatID: "c", Text: "/agent echo2"}
	reply, err := rt.Handle(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(reply, "echo2") {
		t.Fatalf("reply=%q", reply)
	}
	aid, ok := store.GetAgent("stdio:c:u")
	if !ok || aid != "echo2" {
		t.Fatalf("agent=%q ok=%v", aid, ok)
	}
}

func TestToolCommand(t *testing.T) {
	cfg := &config.Config{
		DefaultAgent: "echo",
		Agents:       map[string]config.Agent{"echo": {Mode: "cli", Command: "echo"}},
		Sandbox:      config.Sandbox{AllowCommands: []string{"echo"}},
		Policy:       config.Policy{AllowTools: []string{"echo", "shell"}},
	}
	rt := runtime.New(cfg, runtime.Options{Audit: &policy.BufferAuditor{}})
	out, err := rt.Handle(context.Background(), gateway.Message{
		Channel: "stdio", UserID: "u", ChatID: "c", Text: "/tool echo text=ping",
	})
	if err != nil || out != "ping" {
		t.Fatalf("%v %q", err, out)
	}
}

func TestBindingResolve(t *testing.T) {
	cfg := &config.Config{
		DefaultAgent: "echo",
		Agents: map[string]config.Agent{
			"echo":  {Mode: "cli", Command: "echo"},
			"other": {Mode: "cli", Command: "echo", Args: []string{"other:"}},
		},
		Bindings: []config.Binding{{Channel: "feishu", AgentID: "other"}},
	}
	rt := runtime.New(cfg, runtime.Options{Audit: &policy.BufferAuditor{}})
	out, err := rt.Handle(context.Background(), gateway.Message{
		Channel: "feishu", UserID: "u", ChatID: "c", Text: "hi",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "other:") {
		t.Fatalf("expected other agent, got %q", out)
	}
}

func TestMemoryInjectAndPersistSession(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		DefaultAgent: "echo",
		Agents:       map[string]config.Agent{"echo": {Mode: "cli", Command: "echo"}},
		DataDir:      dir,
		Memory:       config.Memory{Enabled: true, MaxMessages: 10},
	}
	rt, err := runtime.NewFromConfig(cfg, runtime.Options{Audit: &policy.BufferAuditor{}})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if _, err := rt.Handle(ctx, gateway.Message{Channel: "stdio", UserID: "u", ChatID: "c", Text: "first"}); err != nil {
		t.Fatal(err)
	}

	store, err := session.NewFileStore(dir + "/sessions")
	if err != nil {
		t.Fatal(err)
	}
	sid, ok := store.Get("stdio:c:u")
	if !ok || sid == "" {
		t.Fatal("session not persisted")
	}

	reply, err := rt.Handle(ctx, gateway.Message{Channel: "stdio", UserID: "u", ChatID: "c", Text: "second"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(reply, "Prior context:") || !strings.Contains(reply, "first") {
		t.Fatalf("expected memory in prompt echo, got %q", reply)
	}
}
