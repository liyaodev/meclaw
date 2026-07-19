// Package gateway bridges IM channels (WeChat / Feishu / WeCom) to the agent runtime.
package gateway

import "context"

// Message is a normalized inbound IM message.
type Message struct {
	Channel string // wechat | feishu | wecom | stdio | http
	UserID  string
	ChatID  string
	Text    string
	Raw     []byte
}

// Handler processes a normalized inbound message and returns a reply text.
type Handler func(ctx context.Context, msg Message) (string, error)

// Gateway receives IM events and emits replies.
type Gateway interface {
	Start(ctx context.Context) error
	Stop() error
}
