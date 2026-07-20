package feishu_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/meclaw/meclaw/internal/config"
	"github.com/meclaw/meclaw/internal/gateway"
	"github.com/meclaw/meclaw/internal/gateway/feishu"
)

type mockClient struct {
	mu      sync.Mutex
	replies []string
}

func (m *mockClient) ReplyText(ctx context.Context, messageID, text string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.replies = append(m.replies, text)
	return nil
}

func (m *mockClient) SendText(ctx context.Context, chatID, text string) error {
	return nil
}

func TestURLChallenge(t *testing.T) {
	cfg := config.Feishu{VerificationToken: "tok", AppID: "a", AppSecret: "b"}
	adapter := feishu.NewAdapter(cfg, nil, &mockClient{}, nil)
	mux := http.NewServeMux()
	adapter.Mount(mux)

	body, _ := json.Marshal(map[string]string{
		"challenge": "abc",
		"token":     "tok",
		"type":      "url_verification",
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/feishu/event", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d", rr.Code)
	}
	var out map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &out)
	if out["challenge"] != "abc" {
		t.Fatalf("%v", out)
	}
}

func TestInvalidToken(t *testing.T) {
	cfg := config.Feishu{VerificationToken: "tok", AppID: "a", AppSecret: "b"}
	adapter := feishu.NewAdapter(cfg, nil, &mockClient{}, nil)
	mux := http.NewServeMux()
	adapter.Mount(mux)

	body, _ := json.Marshal(map[string]string{
		"challenge": "abc",
		"token":     "wrong",
		"type":      "url_verification",
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/feishu/event", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d", rr.Code)
	}
}

func TestMessageEvent(t *testing.T) {
	cfg := config.Feishu{VerificationToken: "tok", AppID: "a", AppSecret: "b"}
	mc := &mockClient{}
	var got gateway.Message
	handler := func(ctx context.Context, msg gateway.Message) (string, error) {
		got = msg
		return "pong", nil
	}
	adapter := feishu.NewAdapter(cfg, handler, mc, nil)
	mux := http.NewServeMux()
	adapter.Mount(mux)

	event := map[string]any{
		"schema": "2.0",
		"header": map[string]string{
			"event_type": "im.message.receive_v1",
			"token":      "tok",
		},
		"event": map[string]any{
			"sender": map[string]any{
				"sender_id": map[string]string{"open_id": "ou_1"},
			},
			"message": map[string]any{
				"message_id":   "om_1",
				"chat_id":      "oc_1",
				"message_type": "text",
				"content":      `{"text":"@_user_1 hello"}`,
			},
		},
	}
	body, _ := json.Marshal(event)
	req := httptest.NewRequest(http.MethodPost, "/v1/feishu/event", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		mc.mu.Lock()
		n := len(mc.replies)
		mc.mu.Unlock()
		if n > 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	mc.mu.Lock()
	defer mc.mu.Unlock()
	if len(mc.replies) == 0 || mc.replies[0] != "pong" {
		t.Fatalf("replies=%v got=%+v", mc.replies, got)
	}
	if got.UserID != "ou_1" || got.ChatID != "oc_1" || got.Text != "hello" {
		t.Fatalf("msg=%+v", got)
	}
}

func TestParseMessageEvent(t *testing.T) {
	raw := []byte(`{
  "sender":{"sender_id":{"open_id":"ou_x"}},
  "message":{"message_id":"om_x","chat_id":"oc_x","message_type":"text","content":"{\"text\":\"hi\"}"}
}`)
	uid, cid, mid, text, err := feishu.ParseMessageEvent(raw)
	if err != nil {
		t.Fatal(err)
	}
	if uid != "ou_x" || cid != "oc_x" || mid != "om_x" || text != "hi" {
		t.Fatalf("%s %s %s %s", uid, cid, mid, text)
	}
}
