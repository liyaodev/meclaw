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

// HTTPRunner POSTs JSON prompts to a remote agent endpoint.
type HTTPRunner struct {
	client *http.Client
}

// NewHTTPRunner creates an HTTPRunner. If client is nil, a default is used.
func NewHTTPRunner(client *http.Client) *HTTPRunner {
	if client == nil {
		client = &http.Client{Timeout: 120 * time.Second}
	}
	return &HTTPRunner{client: client}
}

// Run is unused directly; use RunURL via Router.
func (r *HTTPRunner) Run(ctx context.Context, req Request) (Response, error) {
	return Response{}, fmt.Errorf("http runner requires base URL via RunURL")
}

type httpRequestBody struct {
	Prompt  string `json:"prompt"`
	Session string `json:"session"`
	AgentID string `json:"agent_id,omitempty"`
}

type httpResponseBody struct {
	Text string `json:"text"`
}

// RunURL POSTs the request to baseURL and returns the text field.
func (r *HTTPRunner) RunURL(ctx context.Context, baseURL string, req Request) (Response, error) {
	if baseURL == "" {
		return Response{}, fmt.Errorf("http: empty base_url")
	}
	body, err := json.Marshal(httpRequestBody{
		Prompt:  req.Prompt,
		Session: req.Session,
		AgentID: req.AgentID,
	})
	if err != nil {
		return Response{}, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(body))
	if err != nil {
		return Response{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(httpReq)
	if err != nil {
		return Response{}, fmt.Errorf("http agent: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Response{}, fmt.Errorf("http agent: status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	var out httpResponseBody
	if err := json.Unmarshal(data, &out); err != nil {
		return Response{}, fmt.Errorf("http agent: decode response: %w", err)
	}
	return Response{Text: out.Text}, nil
}
