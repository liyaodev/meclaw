// Package gateway bridges IM channels (WeChat / Feishu / WeCom) to the agent runtime.
package gateway

// Message is a normalized inbound IM message.
type Message struct {
	Channel   string // wechat | feishu | wecom
	UserID    string
	ChatID    string
	Text      string
	Raw       []byte
}

// Gateway receives IM events and emits replies.
type Gateway interface {
	Start() error
	Stop() error
}
