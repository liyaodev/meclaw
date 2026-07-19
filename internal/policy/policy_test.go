package policy_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/meclaw/meclaw/internal/gateway"
	"github.com/meclaw/meclaw/internal/policy"
)

func TestAllowListEmptyAllowsAll(t *testing.T) {
	a := policy.NewAllowList(nil, nil)
	if !a.AllowMessage(gateway.Message{UserID: "u1"}) {
		t.Fatal("expected allow")
	}
	if !a.AllowTool("u1", "shell") {
		t.Fatal("expected allow tool")
	}
}

func TestAllowListUsers(t *testing.T) {
	a := policy.NewAllowList([]string{"alice"}, nil)
	if !a.AllowMessage(gateway.Message{UserID: "alice"}) {
		t.Fatal("alice should be allowed")
	}
	if a.AllowMessage(gateway.Message{UserID: "bob"}) {
		t.Fatal("bob should be denied")
	}
}

func TestAllowListTools(t *testing.T) {
	a := policy.NewAllowList(nil, []string{"read"})
	if !a.AllowTool("u", "read") {
		t.Fatal("read allowed")
	}
	if a.AllowTool("u", "shell") {
		t.Fatal("shell denied")
	}
}

func TestWriterAuditor(t *testing.T) {
	var buf bytes.Buffer
	a := policy.NewWriterAuditor(&buf)
	a.Log(policy.Event{Action: "allow", UserID: "u1", Detail: "ok"})
	if !strings.Contains(buf.String(), `"action":"allow"`) {
		t.Fatalf("audit=%s", buf.String())
	}
}

func TestBufferAuditor(t *testing.T) {
	a := &policy.BufferAuditor{}
	a.Log(policy.Event{Action: "deny"})
	ev := a.Snapshot()
	if len(ev) != 1 || ev[0].Action != "deny" {
		t.Fatalf("%+v", ev)
	}
}
