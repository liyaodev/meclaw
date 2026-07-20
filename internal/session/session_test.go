package session_test

import (
	"sync"
	"testing"

	"github.com/meclaw/meclaw/internal/session"
)

func TestMemoryStore(t *testing.T) {
	s := session.NewMemoryStore()
	if _, ok := s.Get("k"); ok {
		t.Fatal("expected miss")
	}
	s.Set("k", "s1")
	id, ok := s.Get("k")
	if !ok || id != "s1" {
		t.Fatalf("got %q ok=%v", id, ok)
	}
	s.SetAgent("k", "echo")
	aid, ok := s.GetAgent("k")
	if !ok || aid != "echo" {
		t.Fatalf("agent=%q ok=%v", aid, ok)
	}
	s.Clear("k")
	if _, ok := s.Get("k"); ok {
		t.Fatal("expected cleared")
	}
	if _, ok := s.GetAgent("k"); ok {
		t.Fatal("expected agent cleared")
	}
}

func TestMemoryStoreConcurrent(t *testing.T) {
	s := session.NewMemoryStore()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "k"
			s.Set(key, "sess")
			s.SetAgent(key, "echo")
			_, _ = s.Get(key)
			_, _ = s.GetAgent(key)
		}(i)
	}
	wg.Wait()
}
