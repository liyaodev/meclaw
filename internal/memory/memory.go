// Package memory stores lightweight per-chat conversation memory (scenario A2).
package memory

import "context"

// Message is one memory turn.
type Message struct {
	Role    string `json:"role"` // user | assistant | system
	Content string `json:"content"`
}

// Store appends and recalls chat memory.
type Store interface {
	Append(ctx context.Context, key string, msg Message) error
	Recall(ctx context.Context, key string, limit int) ([]Message, error)
	Clear(ctx context.Context, key string) error
}
