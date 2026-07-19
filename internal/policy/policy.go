// Package policy enforces allow-lists and audit hooks (enterprise / To B).
package policy

import "github.com/meclaw/meclaw/internal/gateway"

// Engine decides whether a message or tool call is allowed.
type Engine interface {
	AllowMessage(msg gateway.Message) bool
	AllowTool(userID, tool string) bool
}

// Auditor records policy and runtime events.
type Auditor interface {
	Log(event Event)
}

// Event is a structured audit record.
type Event struct {
	Action  string `json:"action"` // allow | deny | reply | error
	Channel string `json:"channel,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	ChatID  string `json:"chat_id,omitempty"`
	Agent   string `json:"agent,omitempty"`
	Detail  string `json:"detail,omitempty"`
}
