package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// StateStore is session id + agent selection.
type StateStore interface {
	Store
	AgentStore
}

// FileStore persists session state as JSON under dir/sessions.json.
type FileStore struct {
	mu   sync.RWMutex
	path string
	ids  map[string]string
	ag   map[string]string
}

type filePayload struct {
	IDs    map[string]string `json:"ids"`
	Agents map[string]string `json:"agents"`
}

// NewFileStore loads or creates a file-backed store.
func NewFileStore(dir string) (*FileStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("session dir: %w", err)
	}
	path := filepath.Join(dir, "sessions.json")
	s := &FileStore{
		path: path,
		ids:  map[string]string{},
		ag:   map[string]string{},
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}
	var p filePayload
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parse sessions: %w", err)
	}
	if p.IDs != nil {
		s.ids = p.IDs
	}
	if p.Agents != nil {
		s.ag = p.Agents
	}
	return s, nil
}

func (s *FileStore) persist() error {
	p := filePayload{IDs: s.ids, Agents: s.ag}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// Get returns session id.
func (s *FileStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.ids[key]
	return id, ok
}

// Set stores session id and persists.
func (s *FileStore) Set(key, sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ids[key] = sessionID
	_ = s.persist()
}

// Clear removes key and persists.
func (s *FileStore) Clear(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.ids, key)
	delete(s.ag, key)
	_ = s.persist()
}

// GetAgent returns agent id.
func (s *FileStore) GetAgent(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.ag[key]
	return id, ok
}

// SetAgent stores agent id and persists.
func (s *FileStore) SetAgent(key, agentID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ag[key] = agentID
	_ = s.persist()
}
