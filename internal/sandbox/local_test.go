package sandbox_test

import (
	"context"
	"testing"

	"github.com/meclaw/meclaw/internal/sandbox"
)

func TestLocalExecutorAllow(t *testing.T) {
	ex := sandbox.NewLocalExecutor([]string{"echo"})
	res, err := ex.Exec(context.Background(), sandbox.ExecRequest{Command: "echo", Args: []string{"hi"}})
	if err != nil {
		t.Fatal(err)
	}
	if res.Stdout != "hi" || res.Code != 0 {
		t.Fatalf("%+v", res)
	}
}

func TestLocalExecutorDeny(t *testing.T) {
	ex := sandbox.NewLocalExecutor([]string{"echo"})
	_, err := ex.Exec(context.Background(), sandbox.ExecRequest{Command: "rm", Args: []string{"-rf", "/"}})
	if err == nil {
		t.Fatal("expected deny")
	}
}
