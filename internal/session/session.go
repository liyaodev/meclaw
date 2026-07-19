// Package session tracks per-chat agent conversation state.
package session

// Store maps chat/user keys to agent session ids.
type Store interface {
	Get(key string) (sessionID string, ok bool)
	Set(key, sessionID string)
	Clear(key string)
}

// AgentStore maps chat/user keys to selected agent ids.
type AgentStore interface {
	GetAgent(key string) (agentID string, ok bool)
	SetAgent(key, agentID string)
}
