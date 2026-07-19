// Package agent routes work to ACP / CLI / HTTP agents.
package agent

// Request is a normalized prompt to an agent backend.
type Request struct {
	AgentID string
	Session string
	Prompt  string
	CWD     string
}

// Response is a normalized agent reply.
type Response struct {
	Text string
}

// Runner executes agent requests (ACP preferred over CLI when both exist).
type Runner interface {
	Run(req Request) (Response, error)
}
