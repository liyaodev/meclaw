package orchestrate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/meclaw/meclaw/internal/config"
	"github.com/meclaw/meclaw/internal/orchestrate"
)

func TestRuleResolver(t *testing.T) {
	r := orchestrate.NewRuleResolver([]config.Binding{
		{Channel: "feishu", AgentID: "echo"},
		{Channel: "feishu", UserID: "boss", AgentID: "openai-demo"},
	})
	id, ok := r.Resolve("feishu", "c", "boss")
	if !ok || id != "openai-demo" {
		t.Fatalf("%q %v", id, ok)
	}
	id, ok = r.Resolve("feishu", "c", "other")
	if !ok || id != "echo" {
		t.Fatalf("%q %v", id, ok)
	}
	_, ok = r.Resolve("http", "c", "u")
	if ok {
		t.Fatal("expected miss")
	}
}

func TestSkillLoader(t *testing.T) {
	dir := t.TempDir()
	pack := filepath.Join(dir, "examples", "demo")
	if err := os.MkdirAll(pack, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pack, "README.md"), []byte("skill body"), 0o644); err != nil {
		t.Fatal(err)
	}
	l := orchestrate.SkillLoader{Root: dir}
	text, err := l.Load("examples/demo")
	if err != nil || text != "skill body" {
		t.Fatalf("%v %q", err, text)
	}
}
