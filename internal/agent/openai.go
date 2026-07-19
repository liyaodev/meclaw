package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ChatMessage is an OpenAI-compatible chat turn.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIRunner calls Chat Completions-compatible APIs.
type OpenAIRunner struct {
	client *http.Client
}

// NewOpenAIRunner creates a runner. nil client uses a default timeout.
func NewOpenAIRunner(client *http.Client) *OpenAIRunner {
	if client == nil {
		client = &http.Client{Timeout: 120 * time.Second}
	}
	return &OpenAIRunner{client: client}
}

type openAIReq struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type openAIResp struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// RunChat posts to baseURL (full chat completions URL or API root).
func (r *OpenAIRunner) RunChat(ctx context.Context, baseURL, apiKey, model string, messages []ChatMessage) (Response, error) {
	if baseURL == "" {
		return Response{}, fmt.Errorf("openai: empty base_url")
	}
	if model == "" {
		model = "gpt-4o-mini"
	}
	url := strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(url, "/chat/completions") {
		url += "/chat/completions"
	}
	body, err := json.Marshal(openAIReq{Model: model, Messages: messages})
	if err != nil {
		return Response{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("openai: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}
	var out openAIResp
	if err := json.Unmarshal(data, &out); err != nil {
		return Response{}, fmt.Errorf("openai: decode: %w", err)
	}
	if out.Error != nil && out.Error.Message != "" {
		return Response{}, fmt.Errorf("openai: %s", out.Error.Message)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Response{}, fmt.Errorf("openai: status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	if len(out.Choices) == 0 {
		return Response{}, fmt.Errorf("openai: empty choices")
	}
	return Response{Text: strings.TrimSpace(out.Choices[0].Message.Content)}, nil
}
