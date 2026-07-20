// Package orchestrate binds channels/users to agents and skills (scenario A5).
package orchestrate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/meclaw/meclaw/internal/config"
)

// Binding maps inbound context to an agent id.
type Binding = config.Binding

// Resolver picks an agent for a key.
type Resolver interface {
	Resolve(channel, chatID, userID string) (agentID string, ok bool)
}

// RuleResolver matches the most specific binding (user > chat > channel).
type RuleResolver struct {
	rules []Binding
}

// NewRuleResolver creates a resolver from config bindings.
func NewRuleResolver(rules []Binding) *RuleResolver {
	return &RuleResolver{rules: append([]Binding(nil), rules...)}
}

// Resolve returns the first matching agent id.
func (r *RuleResolver) Resolve(channel, chatID, userID string) (string, bool) {
	bestScore := -1
	best := ""
	for _, b := range r.rules {
		score := 0
		if b.Channel != "" {
			if b.Channel != channel {
				continue
			}
			score++
		}
		if b.ChatID != "" {
			if b.ChatID != chatID {
				continue
			}
			score += 2
		}
		if b.UserID != "" {
			if b.UserID != userID {
				continue
			}
			score += 4
		}
		if score == 0 && b.Channel == "" && b.ChatID == "" && b.UserID == "" {
			continue
		}
		if score > bestScore {
			bestScore = score
			best = b.AgentID
		}
	}
	if bestScore < 0 {
		return "", false
	}
	return best, true
}

// SkillLoader loads skill pack text from disk.
type SkillLoader struct {
	Root string
}

// Load reads README.md (or SKILL.md) under skills_dir/rel.
func (l SkillLoader) Load(rel string) (string, error) {
	if rel == "" {
		return "", fmt.Errorf("empty skill path")
	}
	rel = filepath.Clean(rel)
	if strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("invalid skill path")
	}
	base := filepath.Join(l.Root, rel)
	for _, name := range []string{"SKILL.md", "README.md", "prompt.md"} {
		p := filepath.Join(base, name)
		data, err := os.ReadFile(p)
		if err == nil {
			return string(data), nil
		}
	}
	return "", fmt.Errorf("skill %q: no SKILL.md/README.md/prompt.md", rel)
}

// Attach is a SkillHook that validates the skill exists.
func (l SkillLoader) Attach(agentID, skillPath string) error {
	_, err := l.Load(skillPath)
	if err != nil {
		return fmt.Errorf("attach skill for %s: %w", agentID, err)
	}
	return nil
}
