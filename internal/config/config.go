// Package config loads meclaw runtime configuration.
package config

// Config is the top-level runtime config (scenario A core).
type Config struct {
	DefaultAgent string            `json:"default_agent"`
	Agents       map[string]Agent  `json:"agents"`
	Policy       Policy            `json:"policy"`
}

// Agent describes how to talk to a local or remote agent.
type Agent struct {
	Mode    string   `json:"mode"` // acp | cli | http
	Command string   `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
	BaseURL string   `json:"base_url,omitempty"`
}

// Policy is enterprise-facing access control (scenario A).
type Policy struct {
	AllowTools []string `json:"allow_tools,omitempty"`
	AllowUsers []string `json:"allow_users,omitempty"`
}
