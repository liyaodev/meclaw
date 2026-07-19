// Package eval runs minimal agent regression cases (scenario A4).
package eval

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/meclaw/meclaw/internal/gateway"
)

// Case is one eval prompt/expectation pair.
type Case struct {
	ID       string `json:"id"`
	Prompt   string `json:"prompt"`
	UserID   string `json:"user_id,omitempty"`
	ChatID   string `json:"chat_id,omitempty"`
	Channel  string `json:"channel,omitempty"`
	Contains string `json:"contains"`
}

// Result is a scored case.
type Result struct {
	ID     string `json:"id"`
	Pass   bool   `json:"pass"`
	Detail string `json:"detail"`
}

// Handler is the runtime entry used by eval.
type Handler func(ctx context.Context, msg gateway.Message) (string, error)

// Runner executes eval cases against a handler.
type Runner struct {
	Handle Handler
}

// Run scores all cases.
func (r Runner) Run(ctx context.Context, cases []Case) ([]Result, error) {
	if r.Handle == nil {
		return nil, fmt.Errorf("eval: nil handler")
	}
	out := make([]Result, 0, len(cases))
	for _, c := range cases {
		ch := c.Channel
		if ch == "" {
			ch = "eval"
		}
		uid := c.UserID
		if uid == "" {
			uid = "eval-user"
		}
		cid := c.ChatID
		if cid == "" {
			cid = "eval-chat"
		}
		reply, err := r.Handle(ctx, gateway.Message{
			Channel: ch,
			UserID:  uid,
			ChatID:  cid,
			Text:    c.Prompt,
		})
		res := Result{ID: c.ID}
		if err != nil {
			res.Pass = false
			res.Detail = err.Error()
		} else if c.Contains != "" && !strings.Contains(reply, c.Contains) {
			res.Pass = false
			res.Detail = fmt.Sprintf("reply %q missing %q", truncate(reply, 80), c.Contains)
		} else {
			res.Pass = true
			res.Detail = "ok"
		}
		out = append(out, res)
	}
	return out, nil
}

// LoadCases reads a JSON array of cases from path.
func LoadCases(path string) ([]Case, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cases []Case
	if err := json.Unmarshal(data, &cases); err != nil {
		return nil, fmt.Errorf("parse eval cases: %w", err)
	}
	return cases, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
