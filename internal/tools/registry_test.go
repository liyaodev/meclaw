package tools_test

import (
	"context"
	"testing"

	"github.com/meclaw/meclaw/internal/policy"
	"github.com/meclaw/meclaw/internal/sandbox"
	"github.com/meclaw/meclaw/internal/tools"
)

func TestRegistryEchoAndPolicy(t *testing.T) {
	pol := policy.NewAllowList(nil, []string{"echo"})
	ex := sandbox.NewLocalExecutor([]string{"echo"})
	reg := tools.DefaultRegistry(pol, ex)

	out, err := reg.Call(context.Background(), tools.CallRequest{
		Name: "echo", UserID: "u", Args: map[string]any{"text": "hi"},
	})
	if err != nil || out.Output != "hi" {
		t.Fatalf("%v %+v", err, out)
	}

	_, err = reg.Call(context.Background(), tools.CallRequest{
		Name: "shell", UserID: "u", Args: map[string]any{"command": "echo", "args": "x"},
	})
	if err == nil {
		t.Fatal("shell should be denied by policy")
	}
}

func TestRegistryShell(t *testing.T) {
	pol := policy.NewAllowList(nil, nil)
	ex := sandbox.NewLocalExecutor([]string{"echo"})
	reg := tools.DefaultRegistry(pol, ex)
	out, err := reg.Call(context.Background(), tools.CallRequest{
		Name: "shell", UserID: "u",
		Args: map[string]any{"command": "echo", "args": []string{"ok"}},
	})
	if err != nil || out.Output != "ok" {
		t.Fatalf("%v %+v", err, out)
	}
}
