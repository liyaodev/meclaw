package gateway

import (
	"encoding/json"
	"io"
	"net/http"
)

// HTTPIngress exposes a generic JSON message endpoint.
type HTTPIngress struct {
	Handler Handler
}

type httpMessageRequest struct {
	UserID string `json:"user_id"`
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type httpMessageResponse struct {
	Text  string `json:"text,omitempty"`
	Error string `json:"error,omitempty"`
}

// Mount registers POST /v1/message on mux.
func (h *HTTPIngress) Mount(mux *http.ServeMux) {
	mux.HandleFunc("/v1/message", h.handleMessage)
}

func (h *HTTPIngress) handleMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, httpMessageResponse{Error: "read body"})
		return
	}
	var req httpMessageRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, httpMessageResponse{Error: "invalid json"})
		return
	}
	if req.Text == "" {
		writeJSON(w, http.StatusBadRequest, httpMessageResponse{Error: "text required"})
		return
	}
	if req.UserID == "" {
		req.UserID = "anonymous"
	}
	if req.ChatID == "" {
		req.ChatID = "http"
	}
	reply, err := h.Handler(r.Context(), Message{
		Channel: "http",
		UserID:  req.UserID,
		ChatID:  req.ChatID,
		Text:    req.Text,
		Raw:     body,
	})
	if err != nil {
		writeJSON(w, http.StatusForbidden, httpMessageResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, httpMessageResponse{Text: reply})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
