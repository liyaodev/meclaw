// Package observe provides tracing and health signals (scenario A4).
package observe

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Span is a minimal trace record.
type Span struct {
	Name  string            `json:"name"`
	Start time.Time         `json:"start"`
	End   time.Time         `json:"end"`
	Attrs map[string]string `json:"attrs,omitempty"`
	Error string            `json:"error,omitempty"`
}

// Tracer records spans.
type Tracer interface {
	Start(name string, attrs map[string]string) *Span
	End(span *Span, err error)
}

// NopTracer discards spans.
type NopTracer struct{}

// Start returns an empty span.
func (NopTracer) Start(name string, attrs map[string]string) *Span {
	return &Span{Name: name, Start: time.Now(), Attrs: attrs}
}

// End is a no-op.
func (NopTracer) End(span *Span, err error) {
	_ = span
	_ = err
}

// JSONLTracer appends spans to a JSONL file.
type JSONLTracer struct {
	mu   sync.Mutex
	path string
}

// NewJSONLTracer creates a tracer writing under dir/traces.jsonl.
func NewJSONLTracer(dir string) (*JSONLTracer, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &JSONLTracer{path: filepath.Join(dir, "traces.jsonl")}, nil
}

// Start begins a span.
func (t *JSONLTracer) Start(name string, attrs map[string]string) *Span {
	return &Span{Name: name, Start: time.Now().UTC(), Attrs: attrs}
}

// End writes the span as one JSON line.
func (t *JSONLTracer) End(span *Span, err error) {
	if span == nil {
		return
	}
	span.End = time.Now().UTC()
	if err != nil {
		span.Error = err.Error()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	f, openErr := os.OpenFile(t.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if openErr != nil {
		return
	}
	defer f.Close()
	_ = json.NewEncoder(f).Encode(span)
}

// Health is a JSON health payload.
type Health struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	DataDir string `json:"data_dir,omitempty"`
}

// BuildHealth returns a ready health object.
func BuildHealth(version, dataDir string) Health {
	return Health{Status: "ok", Version: version, DataDir: dataDir}
}
