package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/meclaw/meclaw/internal/config"
)

func TestLoadOK(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.json")
	content := `{
  "default_agent": "echo",
  "agents": {
    "echo": {"mode": "cli", "command": "echo"}
  },
  "policy": {},
  "gateway": {"listen": ":9090"}
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DefaultAgent != "echo" {
		t.Fatalf("default_agent=%q", cfg.DefaultAgent)
	}
	if cfg.Gateway.Listen != ":9090" {
		t.Fatalf("listen=%q", cfg.Gateway.Listen)
	}
}

func TestLoadDefaultListen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.json")
	content := `{
  "default_agent": "echo",
  "agents": {"echo": {"mode": "cli", "command": "echo"}}
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Gateway.Listen != ":8080" {
		t.Fatalf("listen=%q want :8080", cfg.Gateway.Listen)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := config.Load(path); err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadMissingDefaultAgent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.json")
	content := `{"default_agent":"missing","agents":{"echo":{"mode":"cli","command":"echo"}}}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := config.Load(path); err == nil {
		t.Fatal("expected error")
	}
}

func TestFeishuEnabled(t *testing.T) {
	f := config.Feishu{}
	if f.Enabled() {
		t.Fatal("empty should be disabled")
	}
	f = config.Feishu{AppID: "a", AppSecret: "b", VerificationToken: "c"}
	if !f.Enabled() {
		t.Fatal("expected enabled")
	}
}
