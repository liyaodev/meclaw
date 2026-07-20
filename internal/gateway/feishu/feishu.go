// Package feishu adapts Feishu (Lark) bot events to the meclaw gateway.
package feishu

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/meclaw/meclaw/internal/config"
	"github.com/meclaw/meclaw/internal/gateway"
)

// Adapter handles Feishu event callbacks and sends replies.
type Adapter struct {
	cfg     config.Feishu
	handler gateway.Handler
	client  APIClient
	log     *log.Logger
}

// APIClient sends messages via Feishu Open API.
type APIClient interface {
	ReplyText(ctx context.Context, messageID, text string) error
	SendText(ctx context.Context, chatID, text string) error
}

// NewAdapter creates a Feishu adapter.
func NewAdapter(cfg config.Feishu, handler gateway.Handler, client APIClient, logger *log.Logger) *Adapter {
	if logger == nil {
		logger = log.Default()
	}
	if client == nil {
		client = NewClient(cfg, nil)
	}
	return &Adapter{cfg: cfg, handler: handler, client: client, log: logger}
}

// Mount registers POST /v1/feishu/event on mux.
func (a *Adapter) Mount(mux *http.ServeMux) {
	mux.HandleFunc("/v1/feishu/event", a.handleEvent)
}

type challengeReq struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}

type eventEnvelope struct {
	Schema string          `json:"schema"`
	Header eventHeader     `json:"header"`
	Event  json.RawMessage `json:"event"`
	// v1 legacy fields
	Type      string `json:"type"`
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
}

type eventHeader struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	Token      string `json:"token"`
	CreateTime string `json:"create_time"`
}

type messageEvent struct {
	Sender struct {
		SenderID struct {
			OpenID string `json:"open_id"`
			UserID string `json:"user_id"`
		} `json:"sender_id"`
	} `json:"sender"`
	Message struct {
		MessageID   string `json:"message_id"`
		ChatID      string `json:"chat_id"`
		ChatType    string `json:"chat_type"`
		MessageType string `json:"message_type"`
		Content     string `json:"content"`
	} `json:"message"`
}

type textContent struct {
	Text string `json:"text"`
}

func (a *Adapter) handleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}

	// URL verification (challenge) — may arrive as flat JSON.
	var ch challengeReq
	if err := json.Unmarshal(body, &ch); err == nil && ch.Type == "url_verification" {
		if ch.Token != a.cfg.VerificationToken {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"challenge": ch.Challenge})
		return
	}

	var env eventEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	token := env.Header.Token
	if token == "" {
		token = env.Token
	}
	if token != a.cfg.VerificationToken {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	// Acknowledge quickly; process async for message events.
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{}`))

	eventType := env.Header.EventType
	if eventType == "" {
		eventType = env.Type
	}
	if eventType != "im.message.receive_v1" {
		return
	}

	var ev messageEvent
	if err := json.Unmarshal(env.Event, &ev); err != nil {
		a.log.Printf("feishu: parse message event: %v", err)
		return
	}
	if ev.Message.MessageType != "text" {
		return
	}
	var tc textContent
	if err := json.Unmarshal([]byte(ev.Message.Content), &tc); err != nil {
		a.log.Printf("feishu: parse text content: %v", err)
		return
	}
	text := strings.TrimSpace(tc.Text)
	// Strip @bot mentions like @_user_1
	text = stripMention(text)
	if text == "" {
		return
	}

	userID := ev.Sender.SenderID.OpenID
	if userID == "" {
		userID = ev.Sender.SenderID.UserID
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		reply, err := a.handler(ctx, gateway.Message{
			Channel: "feishu",
			UserID:  userID,
			ChatID:  ev.Message.ChatID,
			Text:    text,
			Raw:     body,
		})
		if err != nil {
			a.log.Printf("feishu: handle: %v", err)
			reply = "error: " + err.Error()
		}
		if err := a.client.ReplyText(ctx, ev.Message.MessageID, reply); err != nil {
			a.log.Printf("feishu: reply: %v", err)
			if err2 := a.client.SendText(ctx, ev.Message.ChatID, reply); err2 != nil {
				a.log.Printf("feishu: send: %v", err2)
			}
		}
	}()
}

func stripMention(s string) string {
	fields := strings.Fields(s)
	out := fields[:0]
	for _, f := range fields {
		if strings.HasPrefix(f, "@_") {
			continue
		}
		out = append(out, f)
	}
	return strings.Join(out, " ")
}

// ParseMessageEvent is exported for tests.
func ParseMessageEvent(eventJSON []byte) (userID, chatID, messageID, text string, err error) {
	var ev messageEvent
	if err = json.Unmarshal(eventJSON, &ev); err != nil {
		return
	}
	var tc textContent
	if err = json.Unmarshal([]byte(ev.Message.Content), &tc); err != nil {
		return
	}
	userID = ev.Sender.SenderID.OpenID
	if userID == "" {
		userID = ev.Sender.SenderID.UserID
	}
	return userID, ev.Message.ChatID, ev.Message.MessageID, strings.TrimSpace(tc.Text), nil
}
