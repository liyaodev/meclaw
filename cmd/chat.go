package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/meclaw/meclaw/internal/config"
	"github.com/meclaw/meclaw/internal/gateway"
	"github.com/meclaw/meclaw/internal/policy"
	"github.com/meclaw/meclaw/internal/runtime"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Local stdio chat through the scenario A runtime",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		rt := runtime.New(cfg, runtime.Options{
			Audit: policy.NewWriterAuditor(os.Stderr),
		})
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		stdio := &gateway.Stdio{Handler: rt.Handle}
		if err := stdio.Start(ctx); err != nil && err != context.Canceled {
			return err
		}
		return nil
	},
}

func init() {
	chatCmd.SilenceUsage = true
}
