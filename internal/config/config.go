// Package config loads meclaw runtime configuration.
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config is the top-level runtime config (scenario A1–A6).
type Config struct {
	DefaultAgent string           `json:"default_agent"`
	Agents       map[string]Agent `json:"agents"`
	Policy       Policy           `json:"policy"`
	Gateway      Gateway          `json:"gateway"`
	DataDir      string           `json:"data_dir"`
	Memory       Memory           `json:"memory"`
	Sandbox      Sandbox          `json:"sandbox"`
	Bindings     []Binding        `json:"bindings"`
	SkillsDir    string           `json:"skills_dir"`
}

// Agent describes how to talk to a local or remote agent.
type Agent struct {
	Mode    string   `json:"mode"` // acp | cli | http | openai
	Command string   `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
	BaseURL string   `json:"base_url,omitempty"`
	APIKey  string   `json:"api_key,omitempty"`
	Model   string   `json:"model,omitempty"`
	Skill   string   `json:"skill,omitempty"` // relative path under skills_dir
}

// Memory configures A2 chat memory.
type Memory struct {
	Enabled     bool `json:"enabled"`
	MaxMessages int  `json:"max_messages"`
}

// Sandbox configures A3 local process allow-list.
type Sandbox struct {
	AllowCommands []string `json:"allow_commands"`
}

// Binding maps inbound context to an agent (A5).
type Binding struct {
	Channel string `json:"channel,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	ChatID  string `json:"chat_id,omitempty"`
	AgentID string `json:"agent_id"`
}

// Policy is enterprise-facing access control.
type Policy struct {
	AllowTools []string `json:"allow_tools,omitempty"`
	AllowUsers []string `json:"allow_users,omitempty"`
}

// Gateway holds HTTP listen address and IM channel credentials.
type Gateway struct {
	Listen   string   `json:"listen"`
	Channels Channels `json:"channels"`
}

// Channels groups per-IM adapter settings.
type Channels struct {
	Feishu Feishu `json:"feishu"`
}

// Feishu holds Feishu (Lark) bot credentials.
type Feishu struct {
	AppID             string `json:"app_id"`
	AppSecret         string `json:"app_secret"`
	VerificationToken string `json:"verification_token"`
	EncryptKey        string `json:"encrypt_key,omitempty"`
}

// Enabled reports whether Feishu credentials are present.
func (f Feishu) Enabled() bool {
	return f.AppID != "" && f.AppSecret != "" && f.VerificationToken != ""
}

// Load reads and validates a JSON config file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if cfg.Gateway.Listen == "" {
		cfg.Gateway.Listen = ":8080"
	}
	if cfg.Agents == nil {
		cfg.Agents = map[string]Agent{}
	}
	if cfg.DataDir == "" {
		cfg.DataDir = "./data"
	}
	if cfg.Memory.MaxMessages <= 0 {
		cfg.Memory.MaxMessages = 20
	}
	if cfg.SkillsDir == "" {
		cfg.SkillsDir = "./skills"
	}
	if len(cfg.Sandbox.AllowCommands) == 0 {
		cfg.Sandbox.AllowCommands = []string{"echo", "date", "uname"}
	}
	return &cfg, nil
}

// Validate checks required fields.
func (c *Config) Validate() error {
	if c.DefaultAgent == "" {
		return fmt.Errorf("default_agent is required")
	}
	if len(c.Agents) == 0 {
		return fmt.Errorf("agents must not be empty")
	}
	if _, ok := c.Agents[c.DefaultAgent]; !ok {
		return fmt.Errorf("default_agent %q not found in agents", c.DefaultAgent)
	}
	for id, a := range c.Agents {
		switch a.Mode {
		case "acp", "cli", "http", "openai":
		default:
			return fmt.Errorf("agent %q: unknown mode %q", id, a.Mode)
		}
		if a.Mode == "cli" || a.Mode == "acp" {
			if a.Command == "" {
				return fmt.Errorf("agent %q: command is required for mode %s", id, a.Mode)
			}
		}
		if a.Mode == "http" && a.BaseURL == "" {
			return fmt.Errorf("agent %q: base_url is required for mode http", id)
		}
		if a.Mode == "openai" && a.BaseURL == "" {
			return fmt.Errorf("agent %q: base_url is required for mode openai", id)
		}
	}
	for i, b := range c.Bindings {
		if b.AgentID == "" {
			return fmt.Errorf("bindings[%d]: agent_id required", i)
		}
		if _, ok := c.Agents[b.AgentID]; !ok {
			return fmt.Errorf("bindings[%d]: unknown agent %q", i, b.AgentID)
		}
	}
	return nil
}
