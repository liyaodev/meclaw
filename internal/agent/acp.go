package agent

import (
	"context"
	"fmt"
)

// ACPRunner is a placeholder for Agent Client Protocol backends.
type ACPRunner struct{}

// NewACPRunner creates an ACP stub runner.
func NewACPRunner() *ACPRunner {
	return &ACPRunner{}
}

// Run returns a clear not-implemented error.
func (r *ACPRunner) Run(ctx context.Context, req Request) (Response, error) {
	_ = ctx
	return Response{}, fmt.Errorf("acp runner not implemented (agent %q); use mode cli or http for now", req.AgentID)
}
