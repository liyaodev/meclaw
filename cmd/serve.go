package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/meclaw/meclaw/internal/config"
	"github.com/meclaw/meclaw/internal/gateway"
	"github.com/meclaw/meclaw/internal/gateway/feishu"
	"github.com/meclaw/meclaw/internal/policy"
	"github.com/meclaw/meclaw/internal/runtime"
)

var serveWithStdio bool

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP ingress and optional Feishu webhook",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		rt := runtime.New(cfg, runtime.Options{
			Audit: policy.NewWriterAuditor(os.Stderr),
		})

		mux := http.NewServeMux()
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		ingress := &gateway.HTTPIngress{Handler: rt.Handle}
		ingress.Mount(mux)

		fsCfg := cfg.Gateway.Channels.Feishu
		if fsCfg.Enabled() {
			adapter := feishu.NewAdapter(fsCfg, rt.Handle, nil, log.Default())
			adapter.Mount(mux)
			log.Printf("feishu webhook mounted at /v1/feishu/event")
		} else {
			log.Printf("feishu channel not configured — skipping")
		}

		addr := cfg.Gateway.Listen
		srv := &http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadHeaderTimeout: 10 * time.Second,
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		errCh := make(chan error, 1)
		go func() {
			log.Printf("listening on %s", addr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errCh <- err
			}
			close(errCh)
		}()

		if serveWithStdio {
			go func() {
				stdio := &gateway.Stdio{Handler: rt.Handle}
				_ = stdio.Start(ctx)
				stop()
			}()
		}

		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = srv.Shutdown(shutdownCtx)
			return nil
		case err := <-errCh:
			if err != nil {
				return fmt.Errorf("http server: %w", err)
			}
			return nil
		}
	},
}

func init() {
	serveCmd.Flags().BoolVar(&serveWithStdio, "stdio", false, "also run interactive stdio chat")
	serveCmd.SilenceUsage = true
}
