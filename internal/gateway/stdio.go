package gateway

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

// Stdio runs an interactive chat loop on stdin/stdout.
type Stdio struct {
	Handler Handler
	In      io.Reader
	Out     io.Writer
	UserID  string
	ChatID  string
}

// Start reads lines until EOF or context cancel.
func (s *Stdio) Start(ctx context.Context) error {
	in := s.In
	if in == nil {
		in = os.Stdin
	}
	out := s.Out
	if out == nil {
		out = os.Stdout
	}
	userID := s.UserID
	if userID == "" {
		userID = "local"
	}
	chatID := s.ChatID
	if chatID == "" {
		chatID = "stdio"
	}

	fmt.Fprintln(out, "meclaw chat — type a message (Ctrl-D to quit); /agent <id> to switch")
	sc := bufio.NewScanner(in)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		fmt.Fprint(out, "> ")
		if !sc.Scan() {
			if err := sc.Err(); err != nil {
				return err
			}
			fmt.Fprintln(out)
			return nil
		}
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		reply, err := s.Handler(ctx, Message{
			Channel: "stdio",
			UserID:  userID,
			ChatID:  chatID,
			Text:    line,
		})
		if err != nil {
			fmt.Fprintf(out, "error: %v\n", err)
			continue
		}
		fmt.Fprintln(out, reply)
	}
}

// Stop is a no-op for stdio.
func (s *Stdio) Stop() error { return nil }
