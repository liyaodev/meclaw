package memory_test

import (
	"context"
	"testing"

	"github.com/meclaw/meclaw/internal/memory"
)

func TestFileStoreAppendRecall(t *testing.T) {
	dir := t.TempDir()
	s, err := memory.NewFileStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := s.Append(ctx, "k1", memory.Message{Role: "user", Content: "hi"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Append(ctx, "k1", memory.Message{Role: "assistant", Content: "hello"}); err != nil {
		t.Fatal(err)
	}
	msgs, err := s.Recall(ctx, "k1", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 2 || msgs[0].Content != "hi" || msgs[1].Content != "hello" {
		t.Fatalf("%+v", msgs)
	}
	msgs, err = s.Recall(ctx, "k1", 1)
	if err != nil || len(msgs) != 1 || msgs[0].Content != "hello" {
		t.Fatalf("limit: %+v err=%v", msgs, err)
	}
	if err := s.Clear(ctx, "k1"); err != nil {
		t.Fatal(err)
	}
	msgs, err = s.Recall(ctx, "k1", 0)
	if err != nil || len(msgs) != 0 {
		t.Fatalf("cleared: %+v", msgs)
	}
}
