package agent_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meclaw/meclaw/internal/agent"
	"github.com/meclaw/meclaw/internal/config"
)

func TestCLIRunner(t *testing.T) {
	r := agent.NewCLIRunner()
	resp, err := r.RunCommand(context.Background(), "echo", nil, agent.Request{Prompt: "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text != "hello" {
		t.Fatalf("got %q", resp.Text)
	}
}

func TestHTTPRunner(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		_ = json.NewEncoder(w).Encode(map[string]string{"text": "pong:" + body["prompt"]})
	}))
	defer srv.Close()

	r := agent.NewHTTPRunner(srv.Client())
	resp, err := r.RunURL(context.Background(), srv.URL, agent.Request{Prompt: "hi", Session: "s1"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text != "pong:hi" {
		t.Fatalf("got %q", resp.Text)
	}
}

func TestACPNotImplemented(t *testing.T) {
	r := agent.NewACPRunner()
	_, err := r.Run(context.Background(), agent.Request{AgentID: "claude"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRouterModes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{"text": "from-http"})
	}))
	defer srv.Close()

	router := agent.NewRouter(map[string]config.Agent{
		"echo": {Mode: "cli", Command: "echo"},
		"http": {Mode: "http", BaseURL: srv.URL},
		"acp":  {Mode: "acp", Command: "claude"},
	}).WithHTTPClient(agent.NewHTTPRunner(srv.Client()))

	ctx := context.Background()
	resp, err := router.Run(ctx, agent.Request{AgentID: "echo", Prompt: "x"})
	if err != nil || resp.Text != "x" {
		t.Fatalf("cli: %v %q", err, resp.Text)
	}
	resp, err = router.Run(ctx, agent.Request{AgentID: "http", Prompt: "y"})
	if err != nil || resp.Text != "from-http" {
		t.Fatalf("http: %v %q", err, resp.Text)
	}
	if _, err := router.Run(ctx, agent.Request{AgentID: "acp", Prompt: "z"}); err == nil {
		t.Fatal("acp should fail")
	}
	if _, err := router.Run(ctx, agent.Request{AgentID: "missing", Prompt: "z"}); err == nil {
		t.Fatal("missing should fail")
	}
}
