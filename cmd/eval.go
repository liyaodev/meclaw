package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/meclaw/meclaw/internal/config"
	"github.com/meclaw/meclaw/internal/eval"
	"github.com/meclaw/meclaw/internal/policy"
	"github.com/meclaw/meclaw/internal/runtime"
)

var evalFile string

var evalCmd = &cobra.Command{
	Use:   "eval",
	Short: "Run eval cases against the runtime (scenario A4)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		rt, err := runtime.NewFromConfig(cfg, runtime.Options{
			Audit: policy.NewWriterAuditor(os.Stderr),
		})
		if err != nil {
			return err
		}
		cases, err := eval.LoadCases(evalFile)
		if err != nil {
			return err
		}
		results, err := eval.Runner{Handle: rt.Handle}.Run(context.Background(), cases)
		if err != nil {
			return err
		}
		fail := 0
		for _, r := range results {
			status := "PASS"
			if !r.Pass {
				status = "FAIL"
				fail++
			}
			fmt.Printf("%s %s — %s\n", status, r.ID, r.Detail)
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(results)
		if fail > 0 {
			return fmt.Errorf("%d/%d cases failed", fail, len(results))
		}
		return nil
	},
}

func init() {
	evalCmd.Flags().StringVarP(&evalFile, "file", "f", "evals/smoke.json", "eval cases JSON")
	evalCmd.SilenceUsage = true
}
