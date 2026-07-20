package session

import "sync"

// MemoryStore is a concurrency-safe in-memory session store.
type MemoryStore struct {
	mu     sync.RWMutex
	ids    map[string]string
	agents map[string]string
}

// NewMemoryStore creates an empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		ids:    make(map[string]string),
		agents: make(map[string]string),
	}
}

// Get returns the session id for key.
func (s *MemoryStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.ids[key]
	return id, ok
}

// Set stores the session id for key.
func (s *MemoryStore) Set(key, sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ids[key] = sessionID
}

// Clear removes session and agent selection for key.
func (s *MemoryStore) Clear(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.ids, key)
	delete(s.agents, key)
}

// GetAgent returns the selected agent id for key.
func (s *MemoryStore) GetAgent(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.agents[key]
	return id, ok
}

// SetAgent stores the selected agent id for key.
func (s *MemoryStore) SetAgent(key, agentID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.agents[key] = agentID
}
