package memory

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// FileStore persists memory as JSONL under dir.
type FileStore struct {
	mu  sync.Mutex
	dir string
}

// NewFileStore creates a file-backed memory store.
func NewFileStore(dir string) (*FileStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("memory dir: %w", err)
	}
	return &FileStore{dir: dir}, nil
}

func (s *FileStore) path(key string) string {
	sum := sha256.Sum256([]byte(key))
	return filepath.Join(s.dir, hex.EncodeToString(sum[:16])+".jsonl")
}

// Append writes one message line.
func (s *FileStore) Append(ctx context.Context, key string, msg Message) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.OpenFile(s.path(key), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(msg)
}

// Recall returns the last limit messages (0 = all).
func (s *FileStore) Recall(ctx context.Context, key string, limit int) ([]Message, error) {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.path(key))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var all []Message
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var m Message
		if err := json.Unmarshal(line, &m); err != nil {
			continue
		}
		all = append(all, m)
	}
	if limit > 0 && len(all) > limit {
		all = all[len(all)-limit:]
	}
	return all, nil
}

// Clear removes memory for key.
func (s *FileStore) Clear(ctx context.Context, key string) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	err := os.Remove(s.path(key))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}

// NopStore discards memory.
type NopStore struct{}

// Append is a no-op.
func (NopStore) Append(ctx context.Context, key string, msg Message) error {
	_ = ctx
	_ = key
	_ = msg
	return nil
}

// Recall returns nothing.
func (NopStore) Recall(ctx context.Context, key string, limit int) ([]Message, error) {
	_ = ctx
	_ = key
	_ = limit
	return nil, nil
}

// Clear is a no-op.
func (NopStore) Clear(ctx context.Context, key string) error {
	_ = ctx
	_ = key
	return nil
}
