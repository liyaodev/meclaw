package session_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/meclaw/meclaw/internal/session"
)

func TestFileStorePersist(t *testing.T) {
	dir := t.TempDir()
	s1, err := session.NewFileStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	s1.Set("k", "sess-1")
	s1.SetAgent("k", "echo")

	s2, err := session.NewFileStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	id, ok := s2.Get("k")
	if !ok || id != "sess-1" {
		t.Fatalf("id=%q ok=%v", id, ok)
	}
	aid, ok := s2.GetAgent("k")
	if !ok || aid != "echo" {
		t.Fatalf("agent=%q ok=%v", aid, ok)
	}
	if _, err := os.Stat(filepath.Join(dir, "sessions.json")); err != nil {
		t.Fatal(err)
	}
}
