package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var Version = "dev"

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "meclaw",
	Short: "IM → Agent gateway (Agent Infra)",
	Long: `meclaw is an IM-to-Agent runtime.

Product name: meclaw (claw). Category: Agent Infra / Agent SaaS.
Narrative: IM Agent infrastructure → a work assistant inside WeChat.

Commands:
  chat   local stdio loop
  serve  HTTP ingress + optional Feishu webhook
  version`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "examples/config.example.json", "path to config JSON")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(chatCmd)
	rootCmd.AddCommand(serveCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
