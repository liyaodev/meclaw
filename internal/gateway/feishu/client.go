package feishu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/meclaw/meclaw/internal/config"
)

const defaultBaseURL = "https://open.feishu.cn/open-apis"

// Client talks to Feishu Open API.
type Client struct {
	cfg        config.Feishu
	http       *http.Client
	baseURL    string
	mu         sync.Mutex
	token      string
	tokenExpiry time.Time
}

// NewClient creates a Feishu API client.
func NewClient(cfg config.Feishu, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{cfg: cfg, http: httpClient, baseURL: defaultBaseURL}
}

// WithBaseURL overrides the API base (tests).
func (c *Client) WithBaseURL(u string) *Client {
	c.baseURL = u
	return c
}

type tokenResp struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

func (c *Client) tenantToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" && time.Now().Before(c.tokenExpiry) {
		return c.token, nil
	}
	body, _ := json.Marshal(map[string]string{
		"app_id":     c.cfg.AppID,
		"app_secret": c.cfg.AppSecret,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/auth/v3/tenant_access_token/internal", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var out tokenResp
	if err := json.Unmarshal(data, &out); err != nil {
		return "", err
	}
	if out.Code != 0 || out.TenantAccessToken == "" {
		return "", fmt.Errorf("feishu token: code=%d msg=%s", out.Code, out.Msg)
	}
	c.token = out.TenantAccessToken
	expire := out.Expire
	if expire <= 0 {
		expire = 7200
	}
	c.tokenExpiry = time.Now().Add(time.Duration(expire-60) * time.Second)
	return c.token, nil
}

type apiResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// ReplyText replies to a message by message_id.
func (c *Client) ReplyText(ctx context.Context, messageID, text string) error {
	token, err := c.tenantToken(ctx)
	if err != nil {
		return err
	}
	content, _ := json.Marshal(map[string]string{"text": text})
	payload, _ := json.Marshal(map[string]any{
		"content": string(content),
		"msg_type": "text",
	})
	url := fmt.Sprintf("%s/im/v1/messages/%s/reply", c.baseURL, messageID)
	return c.postJSON(ctx, token, url, payload)
}

// SendText sends a text message to a chat.
func (c *Client) SendText(ctx context.Context, chatID, text string) error {
	token, err := c.tenantToken(ctx)
	if err != nil {
		return err
	}
	content, _ := json.Marshal(map[string]string{"text": text})
	payload, _ := json.Marshal(map[string]any{
		"receive_id": chatID,
		"content":    string(content),
		"msg_type":   "text",
	})
	url := c.baseURL + "/im/v1/messages?receive_id_type=chat_id"
	return c.postJSON(ctx, token, url, payload)
}

func (c *Client) postJSON(ctx context.Context, token, url string, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var out apiResp
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	if out.Code != 0 {
		return fmt.Errorf("feishu api: code=%d msg=%s", out.Code, out.Msg)
	}
	return nil
}
