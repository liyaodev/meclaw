package policy

import "github.com/meclaw/meclaw/internal/gateway"

// AllowList enforces optional user and tool allow-lists.
// Empty lists mean allow-all (developer default).
type AllowList struct {
	Users map[string]struct{}
	Tools map[string]struct{}
}

// NewAllowList builds an AllowList from config slices.
func NewAllowList(users, tools []string) *AllowList {
	a := &AllowList{
		Users: make(map[string]struct{}, len(users)),
		Tools: make(map[string]struct{}, len(tools)),
	}
	for _, u := range users {
		if u != "" {
			a.Users[u] = struct{}{}
		}
	}
	for _, t := range tools {
		if t != "" {
			a.Tools[t] = struct{}{}
		}
	}
	return a
}

// AllowMessage returns true when the user is allowed (or no user list is set).
func (a *AllowList) AllowMessage(msg gateway.Message) bool {
	if len(a.Users) == 0 {
		return true
	}
	_, ok := a.Users[msg.UserID]
	return ok
}

// AllowTool returns true when the tool is allowed (or no tool list is set).
func (a *AllowList) AllowTool(userID, tool string) bool {
	_ = userID
	if len(a.Tools) == 0 {
		return true
	}
	_, ok := a.Tools[tool]
	return ok
}
