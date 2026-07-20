package eval_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/meclaw/meclaw/internal/eval"
	"github.com/meclaw/meclaw/internal/gateway"
)

func TestRunnerPassFail(t *testing.T) {
	r := eval.Runner{Handle: func(ctx context.Context, msg gateway.Message) (string, error) {
		return "meclaw: " + msg.Text, nil
	}}
	res, err := r.Run(context.Background(), []eval.Case{
		{ID: "ok", Prompt: "hi", Contains: "meclaw:"},
		{ID: "bad", Prompt: "hi", Contains: "nope"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res[0].Pass || res[1].Pass {
		t.Fatalf("%+v", res)
	}
}

func TestLoadCases(t *testing.T) {
	path := filepath.Join("..", "..", "evals", "smoke.json")
	cases, err := eval.LoadCases(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("empty")
	}
}
