package policy

import (
	"encoding/json"
	"io"
	"sync"
)

// WriterAuditor writes JSON audit events to an io.Writer.
type WriterAuditor struct {
	mu sync.Mutex
	w  io.Writer
}

// NewWriterAuditor creates an auditor that writes to w.
func NewWriterAuditor(w io.Writer) *WriterAuditor {
	return &WriterAuditor{w: w}
}

// Log writes one audit event as a JSON line.
func (a *WriterAuditor) Log(event Event) {
	a.mu.Lock()
	defer a.mu.Unlock()
	enc := json.NewEncoder(a.w)
	_ = enc.Encode(event)
}

// BufferAuditor stores events in memory for tests.
type BufferAuditor struct {
	mu     sync.Mutex
	Events []Event
}

// Log appends an event.
func (a *BufferAuditor) Log(event Event) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Events = append(a.Events, event)
}

// Snapshot returns a copy of recorded events.
func (a *BufferAuditor) Snapshot() []Event {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]Event, len(a.Events))
	copy(out, a.Events)
	return out
}
