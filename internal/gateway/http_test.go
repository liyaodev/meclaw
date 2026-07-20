package gateway_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meclaw/meclaw/internal/gateway"
)

func TestHTTPIngress(t *testing.T) {
	h := &gateway.HTTPIngress{
		Handler: func(ctx context.Context, msg gateway.Message) (string, error) {
			return "got:" + msg.Text, nil
		},
	}
	mux := http.NewServeMux()
	h.Mount(mux)

	body, _ := json.Marshal(map[string]string{"text": "hi", "user_id": "u"})
	req := httptest.NewRequest(http.MethodPost, "/v1/message", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var out map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &out)
	if out["text"] != "got:hi" {
		t.Fatalf("%v", out)
	}
}
