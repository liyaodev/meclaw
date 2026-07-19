// Package policy enforces allow-lists and audit hooks (enterprise / To B).
package policy

import "github.com/meclaw/meclaw/internal/gateway"

// Engine decides whether a message or tool call is allowed.
type Engine interface {
	AllowMessage(msg gateway.Message) bool
	AllowTool(userID, tool string) bool
}
