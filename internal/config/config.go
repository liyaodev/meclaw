// Package config loads meclaw runtime configuration.
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config is the top-level runtime config (scenario A core).
type Config struct {
	DefaultAgent string           `json:"default_agent"`
	Agents       map[string]Agent `json:"agents"`
	Policy       Policy           `json:"policy"`
	Gateway      Gateway          `json:"gateway"`
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
		case "acp", "cli", "http":
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
	}
	return nil
}
